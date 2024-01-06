// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"os"

	ctrl "sigs.k8s.io/controller-runtime"
	"spheric.cloud/spheric/srictl-machine/cmd/srictl-machine/srictlmachine"
	clicommon "spheric.cloud/spheric/srictl/cmd"
)

func main() {
	ctx := ctrl.SetupSignalHandler()
	if err := srictlmachine.Command(clicommon.OSStreams).ExecuteContext(ctx); err != nil {
		ctrl.Log.Error(err, "Error running command")
		os.Exit(1)
	}
}
