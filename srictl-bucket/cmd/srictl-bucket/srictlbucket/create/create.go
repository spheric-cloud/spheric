// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package create

import (
	"github.com/spf13/cobra"
	"spheric.cloud/spheric/srictl-bucket/cmd/srictl-bucket/srictlbucket/common"
	"spheric.cloud/spheric/srictl-bucket/cmd/srictl-bucket/srictlbucket/create/bucket"
	srictlcmd "spheric.cloud/spheric/srictl/cmd"
)

func Command(streams srictlcmd.Streams, clientFactory common.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use: "create",
	}

	cmd.AddCommand(
		bucket.Command(streams, clientFactory),
	)

	return cmd
}
