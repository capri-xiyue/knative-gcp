/*
Copyright 2019 Google LLC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package pubsub

import (
	"context"
	"fmt"
	"github.com/google/knative-gcp/pkg/reconciler/pubsub"
	"time"

	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"

	"knative.dev/pkg/controller"
	"knative.dev/pkg/logging"
	pkgreconciler "knative.dev/pkg/reconciler"

	"github.com/google/knative-gcp/pkg/apis/events/v1alpha1"
	listers "github.com/google/knative-gcp/pkg/client/listers/events/v1alpha1"
	pubsublisters "github.com/google/knative-gcp/pkg/client/listers/pubsub/v1alpha1"
	"k8s.io/apimachinery/pkg/api/equality"
)

const (
	finalizerName = controllerAgentName

	resourceGroup = "cloudpubsubsources.events.cloud.google.com"
)

// Reconciler is the controller implementation for the CloudPubSubSource source.
type Reconciler struct {
	*pubsub.PubSubBase

	// pubsubLister for reading cloudpubsubsources.
	pubsubLister listers.CloudPubSubSourceLister
	// pullsubscriptionLister for reading pullsubscriptions.
	pullsubscriptionLister pubsublisters.PullSubscriptionLister

}

// Check that we implement the controller.Reconciler interface.
var _ controller.Reconciler = (*Reconciler)(nil)

// Reconcile compares the actual state with the desired, and attempts to
// converge the two. It then updates the Status block of the CloudPubSubSource resource
// with the current status of the resource.
func (r *Reconciler) Reconcile(ctx context.Context, key string) error {
	// Convert the namespace/name string into a distinct namespace and name
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		logging.FromContext(ctx).Desugar().Error("Invalid resource key")
		return nil
	}

	// Get the CloudPubSubSource resource with this namespace/name
	original, err := r.pubsubLister.CloudPubSubSources(namespace).Get(name)
	if apierrs.IsNotFound(err) {
		// The CloudPubSubSource resource may no longer exist, in which case we stop processing.
		logging.FromContext(ctx).Desugar().Error("CloudPubSubSource in work queue no longer exists")
		return nil
	} else if err != nil {
		return err
	}

	// Don't modify the informers copy
	pubsub := original.DeepCopy()

	reconcileErr := r.reconcile(ctx, pubsub)

	// If no error is returned, mark the observed generation.
	if reconcileErr == nil {
		pubsub.Status.ObservedGeneration = pubsub.Generation
	}

	if equality.Semantic.DeepEqual(original.Status, pubsub.Status) {
		// If we didn't change anything then don't call updateStatus.
		// This is important because the copy we loaded from the informer's
		// cache may be stale and we don't want to overwrite a prior update
		// to status with this stale state.

	} else if uErr := r.updateStatus(ctx, original, pubsub); uErr != nil {
		logging.FromContext(ctx).Desugar().Warn("Failed to update CloudPubSubSource status", zap.Error(uErr))
		r.Recorder.Eventf(pubsub, corev1.EventTypeWarning, "UpdateFailed",
			"Failed to update status for CloudPubSubSource %q: %v", pubsub.Name, uErr)
		return uErr
	} else if reconcileErr == nil {
		// There was a difference and updateStatus did not return an error.
		r.Recorder.Eventf(pubsub, corev1.EventTypeNormal, "Updated", "Updated CloudPubSubSource %q", pubsub.Name)
	}
	if reconcileErr != nil {
		r.Recorder.Event(pubsub, corev1.EventTypeWarning, "InternalError", reconcileErr.Error())
	}
	return reconcileErr
}

func (r *Reconciler) reconcile(ctx context.Context, pubsub *v1alpha1.CloudPubSubSource) error {
	ctx = logging.WithLogger(ctx, r.Logger.With(zap.Any("pubsub", pubsub)))

	pubsub.Status.InitializeConditions()

	if pubsub.DeletionTimestamp != nil {
		// No finalizer needed, the pullsubscription will be garbage collected.
		return nil
	}

	_, err := r.PubSubBase.ReconcilePullSubscription(ctx, pubsub, pubsub.Spec.Topic, resourceGroup, true)
	if err != nil {
		logging.FromContext(ctx).Desugar().Error("Failed to reconcile PullSubscription", zap.Error(err))
		return err
	}
	return nil
}


func (r *Reconciler) updateStatus(ctx context.Context, original *v1alpha1.CloudPubSubSource, desired *v1alpha1.CloudPubSubSource) error {
	existing := original.DeepCopy()
	return pkgreconciler.RetryUpdateConflicts(func(attempts int) (err error) {
		// The first iteration tries to use the informer's state, subsequent attempts fetch the latest state via API.
		if attempts > 0 {
			existing, err = r.RunClientSet.EventsV1alpha1().CloudPubSubSources(desired.Namespace).Get(desired.Name, metav1.GetOptions{})
			if err != nil {
				return err
			}
		}
		// Check if there is anything to update.
		if equality.Semantic.DeepEqual(existing.Status, desired.Status) {
			return nil
		}
		becomesReady := desired.Status.IsReady() && !existing.Status.IsReady()

		existing.Status = desired.Status
		_, err = r.RunClientSet.EventsV1alpha1().CloudPubSubSources(desired.Namespace).UpdateStatus(existing)

		if err == nil && becomesReady {
			// TODO compute duration since last non-ready. See https://github.com/google/knative-gcp/issues/455.
			duration := time.Since(existing.ObjectMeta.CreationTimestamp.Time)
			logging.FromContext(ctx).Desugar().Info("CloudPubSubSource became ready", zap.Any("after", duration))
			r.Recorder.Event(existing, corev1.EventTypeNormal, "ReadinessChanged", fmt.Sprintf("CloudPubSubSource %q became ready", existing.Name))
			if metricErr := r.StatsReporter.ReportReady("CloudPubSubSource", existing.Namespace, existing.Name, duration); metricErr != nil {
				logging.FromContext(ctx).Desugar().Error("Failed to record ready for CloudPubSubSource", zap.Error(metricErr))
			}
		}

		return err
	})
}
