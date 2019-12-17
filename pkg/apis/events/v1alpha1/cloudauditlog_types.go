/*
Copyright 2019 Google LLC.

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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"

	duckv1alpha1 "github.com/google/knative-gcp/pkg/apis/duck/v1alpha1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CloudAuditLog is a specification for a Cloud Audit Log event source.
type CloudAuditLog struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CloudAuditLogSpec   `json:"spec"`
	Status CloudAuditLogStatus `json:"status"`
}

var (
	_ runtime.Object = (*CloudAuditLog)(nil)
)

type CloudAuditLogSpec struct {
	// This brings in the PubSub based Source Specs. Includes:
	// Sink, CloudEventOverrides, Secret, PubSubSecret, and Project
	duckv1alpha1.PubSubSpec

	// The CloudAuditLog event source will pull events matching
	// the following parameters:

	// The GCP service providing audit logs. Required.
	Service string
	// The name of the service method or operation. For API calls,
	// this should be the name of the API method. Required.
	Method string
	// The resource or collection that is the target of the
	// operation. The name is a scheme-less URI, not including the
	// API service name.
	Resource string
}

type CloudAuditLogStatus struct {
	duckv1alpha1.PubSubStatus

	SinkID string
}

// GetGroupVersionKind ...
func (_ *CloudAuditLog) GetGroupVersionKind() schema.GroupVersionKind {
	return SchemeGroupVersion.WithKind("CloudAuditLog")
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type CloudAuditLogList struct {
	metav1.TypeMeta
	metav1.ListMeta

	Items []CloudAuditLog `json:"items"`
}
