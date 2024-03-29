// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0

// Code generated by informer-gen. DO NOT EDIT.

package v1alpha1

import (
	internalinterfaces "spheric.cloud/spheric/client-go/informers/internalinterfaces"
)

// Interface provides access to all the informers in this group version.
type Interface interface {
	// ResourceQuotas returns a ResourceQuotaInformer.
	ResourceQuotas() ResourceQuotaInformer
}

type version struct {
	factory          internalinterfaces.SharedInformerFactory
	namespace        string
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// New returns a new Interface.
func New(f internalinterfaces.SharedInformerFactory, namespace string, tweakListOptions internalinterfaces.TweakListOptionsFunc) Interface {
	return &version{factory: f, namespace: namespace, tweakListOptions: tweakListOptions}
}

// ResourceQuotas returns a ResourceQuotaInformer.
func (v *version) ResourceQuotas() ResourceQuotaInformer {
	return &resourceQuotaInformer{factory: v.factory, namespace: v.namespace, tweakListOptions: v.tweakListOptions}
}
