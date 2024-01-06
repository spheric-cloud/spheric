// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package detach

import (
	"github.com/spf13/cobra"
	"spheric.cloud/spheric/srictl-machine/cmd/srictl-machine/srictlmachine/common"
	"spheric.cloud/spheric/srictl-machine/cmd/srictl-machine/srictlmachine/detach/networkinterface"
	"spheric.cloud/spheric/srictl-machine/cmd/srictl-machine/srictlmachine/detach/volume"
	clicommon "spheric.cloud/spheric/srictl/cmd"
)

func Command(streams clicommon.Streams, clientFactory common.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use: "detach",
	}

	cmd.AddCommand(
		volume.Command(streams, clientFactory),
		networkinterface.Command(streams, clientFactory),
	)

	return cmd
}
