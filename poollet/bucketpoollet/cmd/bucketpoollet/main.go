// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"os"

	ctrl "sigs.k8s.io/controller-runtime"
	"spheric.cloud/spheric/poollet/bucketpoollet/cmd/bucketpoollet/app"
)

func main() {
	ctx := ctrl.SetupSignalHandler()
	setupLog := ctrl.Log.WithName("setup")

	if err := app.Command().ExecuteContext(ctx); err != nil {
		setupLog.Error(err, "Error running bucketpoollet")
		os.Exit(1)
	}
}
