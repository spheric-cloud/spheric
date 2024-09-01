// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0

// Code generated by client-gen. DO NOT EDIT.

package v1alpha1

import (
	"context"
	json "encoding/json"
	"fmt"
	"time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
	v1alpha1 "spheric.cloud/spheric/api/core/v1alpha1"
	corev1alpha1 "spheric.cloud/spheric/client-go/applyconfigurations/core/v1alpha1"
	scheme "spheric.cloud/spheric/client-go/spheric/scheme"
)

// FleetsGetter has a method to return a FleetInterface.
// A group's client should implement this interface.
type FleetsGetter interface {
	Fleets() FleetInterface
}

// FleetInterface has methods to work with Fleet resources.
type FleetInterface interface {
	Create(ctx context.Context, fleet *v1alpha1.Fleet, opts v1.CreateOptions) (*v1alpha1.Fleet, error)
	Update(ctx context.Context, fleet *v1alpha1.Fleet, opts v1.UpdateOptions) (*v1alpha1.Fleet, error)
	UpdateStatus(ctx context.Context, fleet *v1alpha1.Fleet, opts v1.UpdateOptions) (*v1alpha1.Fleet, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v1alpha1.Fleet, error)
	List(ctx context.Context, opts v1.ListOptions) (*v1alpha1.FleetList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.Fleet, err error)
	Apply(ctx context.Context, fleet *corev1alpha1.FleetApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.Fleet, err error)
	ApplyStatus(ctx context.Context, fleet *corev1alpha1.FleetApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.Fleet, err error)
	FleetExpansion
}

// fleets implements FleetInterface
type fleets struct {
	client rest.Interface
}

// newFleets returns a Fleets
func newFleets(c *CoreV1alpha1Client) *fleets {
	return &fleets{
		client: c.RESTClient(),
	}
}

// Get takes name of the fleet, and returns the corresponding fleet object, and an error if there is any.
func (c *fleets) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.Fleet, err error) {
	result = &v1alpha1.Fleet{}
	err = c.client.Get().
		Resource("fleets").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of Fleets that match those selectors.
func (c *fleets) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.FleetList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.FleetList{}
	err = c.client.Get().
		Resource("fleets").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested fleets.
func (c *fleets) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Resource("fleets").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a fleet and creates it.  Returns the server's representation of the fleet, and an error, if there is any.
func (c *fleets) Create(ctx context.Context, fleet *v1alpha1.Fleet, opts v1.CreateOptions) (result *v1alpha1.Fleet, err error) {
	result = &v1alpha1.Fleet{}
	err = c.client.Post().
		Resource("fleets").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(fleet).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a fleet and updates it. Returns the server's representation of the fleet, and an error, if there is any.
func (c *fleets) Update(ctx context.Context, fleet *v1alpha1.Fleet, opts v1.UpdateOptions) (result *v1alpha1.Fleet, err error) {
	result = &v1alpha1.Fleet{}
	err = c.client.Put().
		Resource("fleets").
		Name(fleet.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(fleet).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *fleets) UpdateStatus(ctx context.Context, fleet *v1alpha1.Fleet, opts v1.UpdateOptions) (result *v1alpha1.Fleet, err error) {
	result = &v1alpha1.Fleet{}
	err = c.client.Put().
		Resource("fleets").
		Name(fleet.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(fleet).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the fleet and deletes it. Returns an error if one occurs.
func (c *fleets) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Resource("fleets").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *fleets) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Resource("fleets").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched fleet.
func (c *fleets) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.Fleet, err error) {
	result = &v1alpha1.Fleet{}
	err = c.client.Patch(pt).
		Resource("fleets").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}

// Apply takes the given apply declarative configuration, applies it and returns the applied fleet.
func (c *fleets) Apply(ctx context.Context, fleet *corev1alpha1.FleetApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.Fleet, err error) {
	if fleet == nil {
		return nil, fmt.Errorf("fleet provided to Apply must not be nil")
	}
	patchOpts := opts.ToPatchOptions()
	data, err := json.Marshal(fleet)
	if err != nil {
		return nil, err
	}
	name := fleet.Name
	if name == nil {
		return nil, fmt.Errorf("fleet.Name must be provided to Apply")
	}
	result = &v1alpha1.Fleet{}
	err = c.client.Patch(types.ApplyPatchType).
		Resource("fleets").
		Name(*name).
		VersionedParams(&patchOpts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}

// ApplyStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating ApplyStatus().
func (c *fleets) ApplyStatus(ctx context.Context, fleet *corev1alpha1.FleetApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.Fleet, err error) {
	if fleet == nil {
		return nil, fmt.Errorf("fleet provided to Apply must not be nil")
	}
	patchOpts := opts.ToPatchOptions()
	data, err := json.Marshal(fleet)
	if err != nil {
		return nil, err
	}

	name := fleet.Name
	if name == nil {
		return nil, fmt.Errorf("fleet.Name must be provided to Apply")
	}

	result = &v1alpha1.Fleet{}
	err = c.client.Patch(types.ApplyPatchType).
		Resource("fleets").
		Name(*name).
		SubResource("status").
		VersionedParams(&patchOpts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
