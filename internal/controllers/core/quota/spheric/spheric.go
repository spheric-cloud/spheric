// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package spheric

import (
	"spheric.cloud/spheric/internal/controllers/core/quota/compute"
	"spheric.cloud/spheric/internal/controllers/core/quota/generic"
	"spheric.cloud/spheric/internal/controllers/core/quota/storage"
)

var (
	replenishReconcilersBuilder generic.ReplenishReconcilersBuilder
	NewReplenishReconcilers     = replenishReconcilersBuilder.NewReplenishReconcilers
)

func init() {
	replenishReconcilersBuilder.Add(
		compute.NewReplenishReconcilers,
		storage.NewReplenishReconcilers,
	)
}
