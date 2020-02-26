/*
Copyright 2020 Google LLC

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

// Code generated by client-gen. DO NOT EDIT.

package v1alpha1

import (
	"time"

	v1alpha1 "github.com/google/knative-gcp/pkg/apis/events/v1alpha1"
	scheme "github.com/google/knative-gcp/pkg/client/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// CloudBuildSourcesGetter has a method to return a CloudBuildSourceInterface.
// A group's client should implement this interface.
type CloudBuildSourcesGetter interface {
	CloudBuildSources(namespace string) CloudBuildSourceInterface
}

// CloudBuildSourceInterface has methods to work with CloudBuildSource resources.
type CloudBuildSourceInterface interface {
	Create(*v1alpha1.CloudBuildSource) (*v1alpha1.CloudBuildSource, error)
	Update(*v1alpha1.CloudBuildSource) (*v1alpha1.CloudBuildSource, error)
	UpdateStatus(*v1alpha1.CloudBuildSource) (*v1alpha1.CloudBuildSource, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1alpha1.CloudBuildSource, error)
	List(opts v1.ListOptions) (*v1alpha1.CloudBuildSourceList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.CloudBuildSource, err error)
	CloudBuildSourceExpansion
}

// cloudBuildSources implements CloudBuildSourceInterface
type cloudBuildSources struct {
	client rest.Interface
	ns     string
}

// newCloudBuildSources returns a CloudBuildSources
func newCloudBuildSources(c *EventsV1alpha1Client, namespace string) *cloudBuildSources {
	return &cloudBuildSources{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the cloudBuildSource, and returns the corresponding cloudBuildSource object, and an error if there is any.
func (c *cloudBuildSources) Get(name string, options v1.GetOptions) (result *v1alpha1.CloudBuildSource, err error) {
	result = &v1alpha1.CloudBuildSource{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("cloudbuildsources").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of CloudBuildSources that match those selectors.
func (c *cloudBuildSources) List(opts v1.ListOptions) (result *v1alpha1.CloudBuildSourceList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.CloudBuildSourceList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("cloudbuildsources").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested cloudBuildSources.
func (c *cloudBuildSources) Watch(opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("cloudbuildsources").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch()
}

// Create takes the representation of a cloudBuildSource and creates it.  Returns the server's representation of the cloudBuildSource, and an error, if there is any.
func (c *cloudBuildSources) Create(cloudBuildSource *v1alpha1.CloudBuildSource) (result *v1alpha1.CloudBuildSource, err error) {
	result = &v1alpha1.CloudBuildSource{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("cloudbuildsources").
		Body(cloudBuildSource).
		Do().
		Into(result)
	return
}

// Update takes the representation of a cloudBuildSource and updates it. Returns the server's representation of the cloudBuildSource, and an error, if there is any.
func (c *cloudBuildSources) Update(cloudBuildSource *v1alpha1.CloudBuildSource) (result *v1alpha1.CloudBuildSource, err error) {
	result = &v1alpha1.CloudBuildSource{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("cloudbuildsources").
		Name(cloudBuildSource.Name).
		Body(cloudBuildSource).
		Do().
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().

func (c *cloudBuildSources) UpdateStatus(cloudBuildSource *v1alpha1.CloudBuildSource) (result *v1alpha1.CloudBuildSource, err error) {
	result = &v1alpha1.CloudBuildSource{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("cloudbuildsources").
		Name(cloudBuildSource.Name).
		SubResource("status").
		Body(cloudBuildSource).
		Do().
		Into(result)
	return
}

// Delete takes name of the cloudBuildSource and deletes it. Returns an error if one occurs.
func (c *cloudBuildSources) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("cloudbuildsources").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *cloudBuildSources) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	var timeout time.Duration
	if listOptions.TimeoutSeconds != nil {
		timeout = time.Duration(*listOptions.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("cloudbuildsources").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Timeout(timeout).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched cloudBuildSource.
func (c *cloudBuildSources) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.CloudBuildSource, err error) {
	result = &v1alpha1.CloudBuildSource{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("cloudbuildsources").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
