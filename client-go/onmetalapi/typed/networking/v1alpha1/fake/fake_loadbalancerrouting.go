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

package fake

import (
	"context"
	json "encoding/json"
	"fmt"

	v1alpha1 "github.com/onmetal/onmetal-api/api/networking/v1alpha1"
	networkingv1alpha1 "github.com/onmetal/onmetal-api/client-go/applyconfigurations/networking/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeLoadBalancerRoutings implements LoadBalancerRoutingInterface
type FakeLoadBalancerRoutings struct {
	Fake *FakeNetworkingV1alpha1
	ns   string
}

var loadbalancerroutingsResource = v1alpha1.SchemeGroupVersion.WithResource("loadbalancerroutings")

var loadbalancerroutingsKind = v1alpha1.SchemeGroupVersion.WithKind("LoadBalancerRouting")

// Get takes name of the loadBalancerRouting, and returns the corresponding loadBalancerRouting object, and an error if there is any.
func (c *FakeLoadBalancerRoutings) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.LoadBalancerRouting, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(loadbalancerroutingsResource, c.ns, name), &v1alpha1.LoadBalancerRouting{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.LoadBalancerRouting), err
}

// List takes label and field selectors, and returns the list of LoadBalancerRoutings that match those selectors.
func (c *FakeLoadBalancerRoutings) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.LoadBalancerRoutingList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(loadbalancerroutingsResource, loadbalancerroutingsKind, c.ns, opts), &v1alpha1.LoadBalancerRoutingList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.LoadBalancerRoutingList{ListMeta: obj.(*v1alpha1.LoadBalancerRoutingList).ListMeta}
	for _, item := range obj.(*v1alpha1.LoadBalancerRoutingList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested loadBalancerRoutings.
func (c *FakeLoadBalancerRoutings) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(loadbalancerroutingsResource, c.ns, opts))

}

// Create takes the representation of a loadBalancerRouting and creates it.  Returns the server's representation of the loadBalancerRouting, and an error, if there is any.
func (c *FakeLoadBalancerRoutings) Create(ctx context.Context, loadBalancerRouting *v1alpha1.LoadBalancerRouting, opts v1.CreateOptions) (result *v1alpha1.LoadBalancerRouting, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(loadbalancerroutingsResource, c.ns, loadBalancerRouting), &v1alpha1.LoadBalancerRouting{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.LoadBalancerRouting), err
}

// Update takes the representation of a loadBalancerRouting and updates it. Returns the server's representation of the loadBalancerRouting, and an error, if there is any.
func (c *FakeLoadBalancerRoutings) Update(ctx context.Context, loadBalancerRouting *v1alpha1.LoadBalancerRouting, opts v1.UpdateOptions) (result *v1alpha1.LoadBalancerRouting, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(loadbalancerroutingsResource, c.ns, loadBalancerRouting), &v1alpha1.LoadBalancerRouting{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.LoadBalancerRouting), err
}

// Delete takes name of the loadBalancerRouting and deletes it. Returns an error if one occurs.
func (c *FakeLoadBalancerRoutings) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteActionWithOptions(loadbalancerroutingsResource, c.ns, name, opts), &v1alpha1.LoadBalancerRouting{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeLoadBalancerRoutings) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(loadbalancerroutingsResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.LoadBalancerRoutingList{})
	return err
}

// Patch applies the patch and returns the patched loadBalancerRouting.
func (c *FakeLoadBalancerRoutings) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.LoadBalancerRouting, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(loadbalancerroutingsResource, c.ns, name, pt, data, subresources...), &v1alpha1.LoadBalancerRouting{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.LoadBalancerRouting), err
}

// Apply takes the given apply declarative configuration, applies it and returns the applied loadBalancerRouting.
func (c *FakeLoadBalancerRoutings) Apply(ctx context.Context, loadBalancerRouting *networkingv1alpha1.LoadBalancerRoutingApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.LoadBalancerRouting, err error) {
	if loadBalancerRouting == nil {
		return nil, fmt.Errorf("loadBalancerRouting provided to Apply must not be nil")
	}
	data, err := json.Marshal(loadBalancerRouting)
	if err != nil {
		return nil, err
	}
	name := loadBalancerRouting.Name
	if name == nil {
		return nil, fmt.Errorf("loadBalancerRouting.Name must be provided to Apply")
	}
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(loadbalancerroutingsResource, c.ns, *name, types.ApplyPatchType, data), &v1alpha1.LoadBalancerRouting{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.LoadBalancerRouting), err
}
