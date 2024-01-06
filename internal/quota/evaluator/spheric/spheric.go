// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package spheric

import (
	"sigs.k8s.io/controller-runtime/pkg/client"
	"spheric.cloud/spheric/client-go/informers"
	"spheric.cloud/spheric/client-go/spheric"
	"spheric.cloud/spheric/internal/quota/evaluator/compute"
	"spheric.cloud/spheric/internal/quota/evaluator/generic"
	"spheric.cloud/spheric/internal/quota/evaluator/storage"
	"spheric.cloud/spheric/utils/quota"
)

func NewEvaluators(
	machineClassCapabilities,
	volumeClassCapabilities,
	bucketClassCapabilities generic.CapabilitiesReader,
) []quota.Evaluator {
	var evaluators []quota.Evaluator

	evaluators = append(evaluators, compute.NewEvaluators(machineClassCapabilities)...)
	evaluators = append(evaluators, storage.NewEvaluators(volumeClassCapabilities, bucketClassCapabilities)...)

	return evaluators
}

func NewEvaluatorsForAdmission(c spheric.Interface, f informers.SharedInformerFactory) []quota.Evaluator {
	machineClassCapabilities := compute.NewPrimeLRUMachineClassCapabilitiesReader(c, f)
	volumeClassCapabilities := storage.NewPrimeLRUVolumeClassCapabilitiesReader(c, f)
	bucketClassCapabilities := storage.NewPrimeLRUBucketClassCapabilitiesReader(c, f)
	return NewEvaluators(machineClassCapabilities, volumeClassCapabilities, bucketClassCapabilities)
}

func NewEvaluatorsForControllers(c client.Client) []quota.Evaluator {
	machineClassCapabilities := compute.NewClientMachineCapabilitiesReader(c)
	volumeClassCapabilities := storage.NewClientVolumeCapabilitiesReader(c)
	bucketClassCapabilities := storage.NewClientBucketCapabilitiesReader(c)
	return NewEvaluators(machineClassCapabilities, volumeClassCapabilities, bucketClassCapabilities)
}
