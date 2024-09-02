// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package controllers

import (
	"fmt"
	"strconv"

	iri "spheric.cloud/spheric/iri-api/apis/runtime/v1alpha1"
	"spheric.cloud/spheric/utils/generic"
	utilslices "spheric.cloud/spheric/utils/slices"
)

func FindNewIRINetworkInterfaces(desiredIRINics, existingIRINics []*iri.NetworkInterface) []*iri.NetworkInterface {
	var (
		existingIRINicNames = utilslices.ToSetFunc(existingIRINics, (*iri.NetworkInterface).GetName)
		newIRINics          []*iri.NetworkInterface
	)
	for _, desiredIRINic := range desiredIRINics {
		if existingIRINicNames.Has(desiredIRINic.Name) {
			continue
		}

		newIRINics = append(newIRINics, desiredIRINic)
	}
	return newIRINics
}

func FindNewIRIDisks(desiredIRIDisks, existingIRIDisks []*iri.Disk) []*iri.Disk {
	var (
		existingIRIDiskNames = utilslices.ToSetFunc(existingIRIDisks, (*iri.Disk).GetName)
		newIRIDisks          []*iri.Disk
	)
	for _, desiredIRIDisk := range desiredIRIDisks {
		if existingIRIDiskNames.Has(desiredIRIDisk.Name) {
			continue
		}

		newIRIDisks = append(newIRIDisks, desiredIRIDisk)
	}
	return newIRIDisks
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
