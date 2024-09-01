// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and Spheric contributors
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime"
	"spheric.cloud/spheric/api/core/v1alpha1"
)

func addDefaultingFuncs(scheme *runtime.Scheme) error {
	return RegisterDefaults(scheme)
}

func SetDefaults_DiskStatus(status *v1alpha1.DiskStatus) {
	if status.State == "" {
		status.State = v1alpha1.DiskStatePending
	}
}

func SetDefaults_NetworkStatus(status *v1alpha1.NetworkStatus) {
	if status.State == "" {
		status.State = v1alpha1.NetworkStatePending
	}
}

func SetDefaults_NetworkInterfaceStatus(status *v1alpha1.NetworkInterfaceStatus) {
	if status.State == "" {
		status.State = v1alpha1.NetworkInterfaceStatePending
	}
}

func SetDefaults_InstanceStatus(status *v1alpha1.InstanceStatus) {
	if status.State == "" {
		status.State = v1alpha1.InstanceStatePending
	}
}

func SetDefaults_InstanceSpec(spec *v1alpha1.InstanceSpec) {
	if spec.Power == "" {
		spec.Power = v1alpha1.PowerOn
	}
}
