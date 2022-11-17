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

	v1alpha1 "github.com/onmetal/onmetal-api/api/storage/v1alpha1"
	storagev1alpha1 "github.com/onmetal/onmetal-api/client-go/applyconfigurations/storage/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeVolumeClasses implements VolumeClassInterface
type FakeVolumeClasses struct {
	Fake *FakeStorageV1alpha1
}

var volumeclassesResource = schema.GroupVersionResource{Group: "storage.api.onmetal.de", Version: "v1alpha1", Resource: "volumeclasses"}

var volumeclassesKind = schema.GroupVersionKind{Group: "storage.api.onmetal.de", Version: "v1alpha1", Kind: "VolumeClass"}

// Get takes name of the volumeClass, and returns the corresponding volumeClass object, and an error if there is any.
func (c *FakeVolumeClasses) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.VolumeClass, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(volumeclassesResource, name), &v1alpha1.VolumeClass{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.VolumeClass), err
}

// List takes label and field selectors, and returns the list of VolumeClasses that match those selectors.
func (c *FakeVolumeClasses) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.VolumeClassList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(volumeclassesResource, volumeclassesKind, opts), &v1alpha1.VolumeClassList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.VolumeClassList{ListMeta: obj.(*v1alpha1.VolumeClassList).ListMeta}
	for _, item := range obj.(*v1alpha1.VolumeClassList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested volumeClasses.
func (c *FakeVolumeClasses) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(volumeclassesResource, opts))
}

// Create takes the representation of a volumeClass and creates it.  Returns the server's representation of the volumeClass, and an error, if there is any.
func (c *FakeVolumeClasses) Create(ctx context.Context, volumeClass *v1alpha1.VolumeClass, opts v1.CreateOptions) (result *v1alpha1.VolumeClass, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(volumeclassesResource, volumeClass), &v1alpha1.VolumeClass{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.VolumeClass), err
}

// Update takes the representation of a volumeClass and updates it. Returns the server's representation of the volumeClass, and an error, if there is any.
func (c *FakeVolumeClasses) Update(ctx context.Context, volumeClass *v1alpha1.VolumeClass, opts v1.UpdateOptions) (result *v1alpha1.VolumeClass, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(volumeclassesResource, volumeClass), &v1alpha1.VolumeClass{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.VolumeClass), err
}

// Delete takes name of the volumeClass and deletes it. Returns an error if one occurs.
func (c *FakeVolumeClasses) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteActionWithOptions(volumeclassesResource, name, opts), &v1alpha1.VolumeClass{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeVolumeClasses) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(volumeclassesResource, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.VolumeClassList{})
	return err
}

// Patch applies the patch and returns the patched volumeClass.
func (c *FakeVolumeClasses) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.VolumeClass, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(volumeclassesResource, name, pt, data, subresources...), &v1alpha1.VolumeClass{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.VolumeClass), err
}

// Apply takes the given apply declarative configuration, applies it and returns the applied volumeClass.
func (c *FakeVolumeClasses) Apply(ctx context.Context, volumeClass *storagev1alpha1.VolumeClassApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.VolumeClass, err error) {
	if volumeClass == nil {
		return nil, fmt.Errorf("volumeClass provided to Apply must not be nil")
	}
	data, err := json.Marshal(volumeClass)
	if err != nil {
		return nil, err
	}
	name := volumeClass.Name
	if name == nil {
		return nil, fmt.Errorf("volumeClass.Name must be provided to Apply")
	}
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(volumeclassesResource, *name, types.ApplyPatchType, data), &v1alpha1.VolumeClass{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.VolumeClass), err
}
