// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"
	json "encoding/json"
	"fmt"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
	v1alpha1 "spheric.cloud/spheric/api/core/v1alpha1"
	corev1alpha1 "spheric.cloud/spheric/client-go/applyconfigurations/core/v1alpha1"
)

// FakeFleets implements FleetInterface
type FakeFleets struct {
	Fake *FakeCoreV1alpha1
}

var fleetsResource = v1alpha1.SchemeGroupVersion.WithResource("fleets")

var fleetsKind = v1alpha1.SchemeGroupVersion.WithKind("Fleet")

// Get takes name of the fleet, and returns the corresponding fleet object, and an error if there is any.
func (c *FakeFleets) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.Fleet, err error) {
	emptyResult := &v1alpha1.Fleet{}
	obj, err := c.Fake.
		Invokes(testing.NewRootGetActionWithOptions(fleetsResource, name, options), emptyResult)
	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1alpha1.Fleet), err
}

// List takes label and field selectors, and returns the list of Fleets that match those selectors.
func (c *FakeFleets) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.FleetList, err error) {
	emptyResult := &v1alpha1.FleetList{}
	obj, err := c.Fake.
		Invokes(testing.NewRootListActionWithOptions(fleetsResource, fleetsKind, opts), emptyResult)
	if obj == nil {
		return emptyResult, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.FleetList{ListMeta: obj.(*v1alpha1.FleetList).ListMeta}
	for _, item := range obj.(*v1alpha1.FleetList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested fleets.
func (c *FakeFleets) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchActionWithOptions(fleetsResource, opts))
}

// Create takes the representation of a fleet and creates it.  Returns the server's representation of the fleet, and an error, if there is any.
func (c *FakeFleets) Create(ctx context.Context, fleet *v1alpha1.Fleet, opts v1.CreateOptions) (result *v1alpha1.Fleet, err error) {
	emptyResult := &v1alpha1.Fleet{}
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateActionWithOptions(fleetsResource, fleet, opts), emptyResult)
	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1alpha1.Fleet), err
}

// Update takes the representation of a fleet and updates it. Returns the server's representation of the fleet, and an error, if there is any.
func (c *FakeFleets) Update(ctx context.Context, fleet *v1alpha1.Fleet, opts v1.UpdateOptions) (result *v1alpha1.Fleet, err error) {
	emptyResult := &v1alpha1.Fleet{}
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateActionWithOptions(fleetsResource, fleet, opts), emptyResult)
	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1alpha1.Fleet), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeFleets) UpdateStatus(ctx context.Context, fleet *v1alpha1.Fleet, opts v1.UpdateOptions) (result *v1alpha1.Fleet, err error) {
	emptyResult := &v1alpha1.Fleet{}
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateSubresourceActionWithOptions(fleetsResource, "status", fleet, opts), emptyResult)
	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1alpha1.Fleet), err
}

// Delete takes name of the fleet and deletes it. Returns an error if one occurs.
func (c *FakeFleets) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteActionWithOptions(fleetsResource, name, opts), &v1alpha1.Fleet{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeFleets) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionActionWithOptions(fleetsResource, opts, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.FleetList{})
	return err
}

// Patch applies the patch and returns the patched fleet.
func (c *FakeFleets) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.Fleet, err error) {
	emptyResult := &v1alpha1.Fleet{}
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceActionWithOptions(fleetsResource, name, pt, data, opts, subresources...), emptyResult)
	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1alpha1.Fleet), err
}

// Apply takes the given apply declarative configuration, applies it and returns the applied fleet.
func (c *FakeFleets) Apply(ctx context.Context, fleet *corev1alpha1.FleetApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.Fleet, err error) {
	if fleet == nil {
		return nil, fmt.Errorf("fleet provided to Apply must not be nil")
	}
	data, err := json.Marshal(fleet)
	if err != nil {
		return nil, err
	}
	name := fleet.Name
	if name == nil {
		return nil, fmt.Errorf("fleet.Name must be provided to Apply")
	}
	emptyResult := &v1alpha1.Fleet{}
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceActionWithOptions(fleetsResource, *name, types.ApplyPatchType, data, opts.ToPatchOptions()), emptyResult)
	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1alpha1.Fleet), err
}

// ApplyStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating ApplyStatus().
func (c *FakeFleets) ApplyStatus(ctx context.Context, fleet *corev1alpha1.FleetApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.Fleet, err error) {
	if fleet == nil {
		return nil, fmt.Errorf("fleet provided to Apply must not be nil")
	}
	data, err := json.Marshal(fleet)
	if err != nil {
		return nil, err
	}
	name := fleet.Name
	if name == nil {
		return nil, fmt.Errorf("fleet.Name must be provided to Apply")
	}
	emptyResult := &v1alpha1.Fleet{}
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceActionWithOptions(fleetsResource, *name, types.ApplyPatchType, data, opts.ToPatchOptions(), "status"), emptyResult)
	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1alpha1.Fleet), err
}
