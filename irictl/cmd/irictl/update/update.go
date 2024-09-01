// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package update

import (
	"github.com/spf13/cobra"
	clicommon "spheric.cloud/spheric/irictl/cmd"
	"spheric.cloud/spheric/irictl/cmd/irictl/common"
)

func Command(streams clicommon.Streams, clientFactory common.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use: "update",
	}

	cmd.AddCommand()

	return cmd
}