// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package controllers

import (
	"fmt"
	"strconv"

	sri "spheric.cloud/spheric/sri/apis/machine/v1alpha1"
	"spheric.cloud/spheric/utils/generic"
	utilslices "spheric.cloud/spheric/utils/slices"
)

func FindNewSRINetworkInterfaces(desiredSRINics, existingSRINics []*sri.NetworkInterface) []*sri.NetworkInterface {
	var (
		existingSRINicNames = utilslices.ToSetFunc(existingSRINics, (*sri.NetworkInterface).GetName)
		newSRINics          []*sri.NetworkInterface
	)
	for _, desiredSRINic := range desiredSRINics {
		if existingSRINicNames.Has(desiredSRINic.Name) {
			continue
		}

		newSRINics = append(newSRINics, desiredSRINic)
	}
	return newSRINics
}

func FindNewSRIVolumes(desiredSRIVolumes, existingSRIVolumes []*sri.Volume) []*sri.Volume {
	var (
		existingSRIVolumeNames = utilslices.ToSetFunc(existingSRIVolumes, (*sri.Volume).GetName)
		newSRIVolumes          []*sri.Volume
	)
	for _, desiredSRIVolume := range desiredSRIVolumes {
		if existingSRIVolumeNames.Has(desiredSRIVolume.Name) {
			continue
		}

		newSRIVolumes = append(newSRIVolumes, desiredSRIVolume)
	}
	return newSRIVolumes
}

func parseInt64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

func getAndParseFromStringMap[E any](annotations map[string]string, key string, parse func(string) (E, error)) (E, error) {
	s, ok := annotations[key]
	if !ok {
		return generic.Zero[E](), fmt.Errorf("no value found at key %s", key)
	}

	e, err := parse(s)
	if err != nil {
		return e, fmt.Errorf("error parsing key %s data %s: %w", key, s, err)
	}

	return e, nil
}
