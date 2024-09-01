// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package initializer

import (
	"k8s.io/apiserver/pkg/admission"
	sphericinformers "spheric.cloud/spheric/client-go/informers"
	"spheric.cloud/spheric/client-go/spheric"
)

type initializer struct {
	externalClient    spheric.Interface
	externalInformers sphericinformers.SharedInformerFactory
}

func New(
	externalClient spheric.Interface,
	externalInformers sphericinformers.SharedInformerFactory,
) admission.PluginInitializer {
	return &initializer{
		externalClient:    externalClient,
		externalInformers: externalInformers,
	}
}

func (i *initializer) Initialize(plugin admission.Interface) {
	if wants, ok := plugin.(WantsExternalSphericClientSet); ok {
		wants.SetExternalSphericClientSet(i.externalClient)
	}

	if wants, ok := plugin.(WantsExternalInformers); ok {
		wants.SetExternalSphericInformerFactory(i.externalInformers)
	}
}

type WantsExternalSphericClientSet interface {
	SetExternalSphericClientSet(client spheric.Interface)
	admission.InitializationValidator
}

type WantsExternalInformers interface {
	SetExternalSphericInformerFactory(f sphericinformers.SharedInformerFactory)
	admission.InitializationValidator
}
