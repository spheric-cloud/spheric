// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package tableconverters

import (
	"spheric.cloud/spheric/irictl/tableconverter"
)

var (
	RegistryBuilder tableconverter.RegistryBuilder
	AddToRegistry   = RegistryBuilder.AddToRegistry
)
