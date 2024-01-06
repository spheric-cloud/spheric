// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"os"

	ctrl "sigs.k8s.io/controller-runtime"
	"spheric.cloud/spheric/srictl-bucket/cmd/srictl-bucket/srictlbucket"
	srictlcmd "spheric.cloud/spheric/srictl/cmd"
)

func main() {
	ctx := ctrl.SetupSignalHandler()
	if err := srictlbucket.Command(srictlcmd.OSStreams).ExecuteContext(ctx); err != nil {
		ctrl.Log.Error(err, "Error running command")
		os.Exit(1)
	}
}
