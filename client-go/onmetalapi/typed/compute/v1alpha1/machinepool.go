/*
 * Copyright (c) 2022 by the OnMetal authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
// Code generated by client-gen. DO NOT EDIT.

package v1alpha1

import (
	"context"
	json "encoding/json"
	"fmt"
	"time"

	v1alpha1 "github.com/onmetal/onmetal-api/api/compute/v1alpha1"
	computev1alpha1 "github.com/onmetal/onmetal-api/client-go/applyconfigurations/compute/v1alpha1"
	scheme "github.com/onmetal/onmetal-api/client-go/onmetalapi/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// MachinePoolsGetter has a method to return a MachinePoolInterface.
// A group's client should implement this interface.
type MachinePoolsGetter interface {
	MachinePools() MachinePoolInterface
}

// MachinePoolInterface has methods to work with MachinePool resources.
type MachinePoolInterface interface {
	Create(ctx context.Context, machinePool *v1alpha1.MachinePool, opts v1.CreateOptions) (*v1alpha1.MachinePool, error)
	Update(ctx context.Context, machinePool *v1alpha1.MachinePool, opts v1.UpdateOptions) (*v1alpha1.MachinePool, error)
	UpdateStatus(ctx context.Context, machinePool *v1alpha1.MachinePool, opts v1.UpdateOptions) (*v1alpha1.MachinePool, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v1alpha1.MachinePool, error)
	List(ctx context.Context, opts v1.ListOptions) (*v1alpha1.MachinePoolList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.MachinePool, err error)
	Apply(ctx context.Context, machinePool *computev1alpha1.MachinePoolApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.MachinePool, err error)
	ApplyStatus(ctx context.Context, machinePool *computev1alpha1.MachinePoolApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.MachinePool, err error)
	MachinePoolExpansion
}

// machinePools implements MachinePoolInterface
type machinePools struct {
	client rest.Interface
}

// newMachinePools returns a MachinePools
func newMachinePools(c *ComputeV1alpha1Client) *machinePools {
	return &machinePools{
		client: c.RESTClient(),
	}
}

// Get takes name of the machinePool, and returns the corresponding machinePool object, and an error if there is any.
func (c *machinePools) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.MachinePool, err error) {
	result = &v1alpha1.MachinePool{}
	err = c.client.Get().
		Resource("machinepools").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of MachinePools that match those selectors.
func (c *machinePools) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.MachinePoolList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.MachinePoolList{}
	err = c.client.Get().
		Resource("machinepools").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested machinePools.
func (c *machinePools) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Resource("machinepools").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a machinePool and creates it.  Returns the server's representation of the machinePool, and an error, if there is any.
func (c *machinePools) Create(ctx context.Context, machinePool *v1alpha1.MachinePool, opts v1.CreateOptions) (result *v1alpha1.MachinePool, err error) {
	result = &v1alpha1.MachinePool{}
	err = c.client.Post().
		Resource("machinepools").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(machinePool).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a machinePool and updates it. Returns the server's representation of the machinePool, and an error, if there is any.
func (c *machinePools) Update(ctx context.Context, machinePool *v1alpha1.MachinePool, opts v1.UpdateOptions) (result *v1alpha1.MachinePool, err error) {
	result = &v1alpha1.MachinePool{}
	err = c.client.Put().
		Resource("machinepools").
		Name(machinePool.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(machinePool).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *machinePools) UpdateStatus(ctx context.Context, machinePool *v1alpha1.MachinePool, opts v1.UpdateOptions) (result *v1alpha1.MachinePool, err error) {
	result = &v1alpha1.MachinePool{}
	err = c.client.Put().
		Resource("machinepools").
		Name(machinePool.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(machinePool).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the machinePool and deletes it. Returns an error if one occurs.
func (c *machinePools) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Resource("machinepools").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *machinePools) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Resource("machinepools").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched machinePool.
func (c *machinePools) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.MachinePool, err error) {
	result = &v1alpha1.MachinePool{}
	err = c.client.Patch(pt).
		Resource("machinepools").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}

// Apply takes the given apply declarative configuration, applies it and returns the applied machinePool.
func (c *machinePools) Apply(ctx context.Context, machinePool *computev1alpha1.MachinePoolApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.MachinePool, err error) {
	if machinePool == nil {
		return nil, fmt.Errorf("machinePool provided to Apply must not be nil")
	}
	patchOpts := opts.ToPatchOptions()
	data, err := json.Marshal(machinePool)
	if err != nil {
		return nil, err
	}
	name := machinePool.Name
	if name == nil {
		return nil, fmt.Errorf("machinePool.Name must be provided to Apply")
	}
	result = &v1alpha1.MachinePool{}
	err = c.client.Patch(types.ApplyPatchType).
		Resource("machinepools").
		Name(*name).
		VersionedParams(&patchOpts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}

// ApplyStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating ApplyStatus().
func (c *machinePools) ApplyStatus(ctx context.Context, machinePool *computev1alpha1.MachinePoolApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.MachinePool, err error) {
	if machinePool == nil {
		return nil, fmt.Errorf("machinePool provided to Apply must not be nil")
	}
	patchOpts := opts.ToPatchOptions()
	data, err := json.Marshal(machinePool)
	if err != nil {
		return nil, err
	}

	name := machinePool.Name
	if name == nil {
		return nil, fmt.Errorf("machinePool.Name must be provided to Apply")
	}

	result = &v1alpha1.MachinePool{}
	err = c.client.Patch(types.ApplyPatchType).
		Resource("machinepools").
		Name(*name).
		SubResource("status").
		VersionedParams(&patchOpts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
