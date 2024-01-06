// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package compute

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	computev1alpha1 "spheric.cloud/spheric/api/compute/v1alpha1"
	corev1alpha1 "spheric.cloud/spheric/api/core/v1alpha1"
	"spheric.cloud/spheric/client-go/informers"
	"spheric.cloud/spheric/client-go/spheric"
	"spheric.cloud/spheric/internal/quota/evaluator/generic"
	utilsgeneric "spheric.cloud/spheric/utils/generic"
	"spheric.cloud/spheric/utils/quota"
	"spheric.cloud/spheric/utils/quota/resourceaccess"
)

func NewEvaluators(capabilities generic.CapabilitiesReader) []quota.Evaluator {
	return []quota.Evaluator{
		NewMachineEvaluator(capabilities),
	}
}

func extractMachineClassCapabilities(machineClass *computev1alpha1.MachineClass) corev1alpha1.ResourceList {
	return machineClass.Capabilities
}

func NewClientMachineCapabilitiesReader(c client.Client) generic.CapabilitiesReader {
	getter := resourceaccess.NewTypedClientGetter[computev1alpha1.MachineClass](c)
	return generic.NewGetterCapabilitiesReader(getter,
		extractMachineClassCapabilities,
		func(s string) client.ObjectKey { return client.ObjectKey{Name: s} },
	)
}

func NewPrimeLRUMachineClassCapabilitiesReader(c spheric.Interface, f informers.SharedInformerFactory) generic.CapabilitiesReader {
	getter := resourceaccess.NewPrimeLRUGetter[*computev1alpha1.MachineClass, string](
		func(ctx context.Context, className string) (*computev1alpha1.MachineClass, error) {
			return c.ComputeV1alpha1().MachineClasses().Get(ctx, className, metav1.GetOptions{})
		},
		func(ctx context.Context, className string) (*computev1alpha1.MachineClass, error) {
			return f.Compute().V1alpha1().MachineClasses().Lister().Get(className)
		},
	)
	return generic.NewGetterCapabilitiesReader(getter, extractMachineClassCapabilities, utilsgeneric.Identity[string])
}
