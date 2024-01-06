// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package status

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	ctrl "sigs.k8s.io/controller-runtime"
	sri "spheric.cloud/spheric/sri/apis/machine/v1alpha1"
	"spheric.cloud/spheric/srictl-machine/cmd/srictl-machine/srictlmachine/common"
	clicommon "spheric.cloud/spheric/srictl/cmd"
	"spheric.cloud/spheric/srictl/renderer"
)

func Command(streams clicommon.Streams, clientFactory common.Factory) *cobra.Command {
	var (
		outputOpts = clientFactory.OutputOptions()
	)

	cmd := &cobra.Command{
		Use: "status",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			log := ctrl.LoggerFrom(ctx)

			client, cleanup, err := clientFactory.Client()
			if err != nil {
				return err
			}
			defer func() {
				if err := cleanup(); err != nil {
					log.Error(err, "Error cleaning up")
				}
			}()

			render, err := outputOpts.Renderer("table")
			if err != nil {
				return err
			}

			return Run(cmd.Context(), streams, client, render)
		},
	}

	outputOpts.AddFlags(cmd.Flags())

	return cmd
}

func Run(ctx context.Context, streams clicommon.Streams, client sri.MachineRuntimeClient, render renderer.Renderer) error {
	res, err := client.Status(ctx, &sri.StatusRequest{})
	if err != nil {
		return fmt.Errorf("error getting status: %w", err)
	}

	return render.Render(res.MachineClassStatus, streams.Out)
}
