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

package v1

import (
	duckv1 "github.com/google/knative-gcp/pkg/apis/duck/v1"
)

// PubSubable is an interface that each duckv1.PubSub duck type must
// support in order to get reconciled properly in a generic way.
type PubSubable interface {
	Identifiable
	// PubSubSpec returns the PubSubSpec portion of the Spec.
	PubSubSpec() *duckv1.PubSubSpec
	// PubSubStatus returns the PubSubStatus portion of the Status.
	PubSubStatus() *duckv1.PubSubStatus
}
