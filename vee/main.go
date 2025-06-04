// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"os"

	ctrl "sigs.k8s.io/controller-runtime"
	"spheric.cloud/spheric/vee/cmd/vee"
)

func main() {
	ctx := ctrl.SetupSignalHandler()
	log := ctrl.Log.WithName("main")

	if err := vee.Command().ExecuteContext(ctx); err != nil {
		log.V(1).Error(err, "Error running vee")
		os.Exit(1)
	}
}
