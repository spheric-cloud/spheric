// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package instancediskdevices

import (
	"context"
	"fmt"
	"io"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apiserver/pkg/admission"
	"spheric.cloud/spheric/internal/admission/plugin/instancediskdevices/device"
	"spheric.cloud/spheric/internal/apis/core"
)

// PluginName indicates name of admission plugin.
const PluginName = "InstanceDiskDevices"

// Register registers a plugin
func Register(plugins *admission.Plugins) {
	plugins.Register(PluginName, func(config io.Reader) (admission.Interface, error) {
		return NewInstanceDiskDevices(), nil
	})
}

type InstanceDiskDevices struct {
	*admission.Handler
}

func NewInstanceDiskDevices() *InstanceDiskDevices {
	return &InstanceDiskDevices{
		Handler: admission.NewHandler(admission.Create, admission.Update),
	}
}

func (d *InstanceDiskDevices) Admit(ctx context.Context, a admission.Attributes, o admission.ObjectInterfaces) error {
	if shouldIgnore(a) {
		return nil
	}

	instance, ok := a.GetObject().(*core.Instance)
	if !ok {
		return apierrors.NewBadRequest("Resource was marked with kind Instance but was unable to be converted")
	}

	namer, err := deviceNamerFromInstanceDisks(instance)
	if err != nil {
		return apierrors.NewBadRequest("Instance has conflicting disk device names")
	}

	for i := range instance.Spec.Disks {
		disk := &instance.Spec.Disks[i]
		if disk.Device != nil && *disk.Device != "" {
			continue
		}

		newDevice, err := namer.Generate(device.SphericPrefix) // TODO: We should have a better way for a device prefix.
		if err != nil {
			return apierrors.NewBadRequest("No device names left for instance")
		}

		disk.Device = &newDevice
	}

	return nil
}

func shouldIgnore(a admission.Attributes) bool {
	if a.GetKind().GroupKind() != core.Kind("Instance").GroupKind() {
		return true
	}

	instance, ok := a.GetObject().(*core.Instance)
	if !ok {
		return true
	}

	return !instanceHasAnyDiskWithoutDevice(instance)
}

func instanceHasAnyDiskWithoutDevice(instance *core.Instance) bool {
	for _, disk := range instance.Spec.Disks {
		if disk.Device == nil || *disk.Device == "" {
			return true
		}
	}
	return false
}

func deviceNamerFromInstanceDisks(instance *core.Instance) (*device.Namer, error) {
	namer := device.NewNamer()
	for _, disk := range instance.Spec.Disks {
		if dev := disk.Device; dev != nil && *dev != "" {
			if err := namer.Observe(*dev); err != nil {
				return nil, fmt.Errorf("error observing device %s: %w", *dev, err)
			}
		}
	}
	return namer, nil
}
