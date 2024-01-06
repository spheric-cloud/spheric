// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package compute

import (
	computev1alpha1 "spheric.cloud/spheric/api/compute/v1alpha1"
	"spheric.cloud/spheric/internal/controllers/core/quota/generic"
)

var (
	replenishReconcilersBuilder generic.ReplenishReconcilersBuilder
	NewReplenishReconcilers     = replenishReconcilersBuilder.NewReplenishReconcilers
)

func init() {
	replenishReconcilersBuilder.Register(
		&computev1alpha1.Machine{},
	)
}
