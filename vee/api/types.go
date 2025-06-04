// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0

package api

import "spheric.cloud/spheric/actuo/meta"

type Instance struct {
	meta.ObjectMeta `json:"metadata,omitempty"`
	ID              string       `json:"id"`
	Spec            InstanceSpec `json:"spec"`
}

type InstanceSpec struct {
	CPUCount    int32 `json:"cpuCount"`
	MemoryBytes int64 `json:"memoryBytes"`
}
