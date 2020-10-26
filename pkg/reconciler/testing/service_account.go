/*
Copyright 2020 The Knative Authors

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

package testing

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ServiceAccountOption func(*corev1.ServiceAccount)

// NewServiceAccount creates a ServiceAccount
func NewServiceAccount(name, namespace string, so ...ServiceAccountOption) *corev1.ServiceAccount {
	sa := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   namespace,
			Annotations: map[string]string{},
		},
	}
	for _, opt := range so {
		opt(sa)
	}
	return sa
}

func WithServiceAccountAnnotation(gServiceAccount string) ServiceAccountOption {
	return func(s *corev1.ServiceAccount) {
		s.Annotations["iam.gke.io/gcp-service-account"] = gServiceAccount
	}
}

func WithServiceAccountOwnerReferences(ownerReferences []metav1.OwnerReference) ServiceAccountOption {
	return func(s *corev1.ServiceAccount) {
		s.OwnerReferences = ownerReferences
	}
}
