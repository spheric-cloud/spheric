// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package machine

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	ctrl "sigs.k8s.io/controller-runtime"
	sri "spheric.cloud/spheric/sri/apis/machine/v1alpha1"
	"spheric.cloud/spheric/srictl-machine/cmd/srictl-machine/srictlmachine/common"
	clicommon "spheric.cloud/spheric/srictl/cmd"
)

func Command(streams clicommon.Streams, clientFactory common.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "machine id [ids...]",
		Aliases: common.MachineAliases,
		Args:    cobra.MinimumNArgs(1),
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

			ids := args

			return Run(cmd.Context(), streams, client, ids)
		},
	}

	return cmd
}

func Run(ctx context.Context, streams clicommon.Streams, client sri.MachineRuntimeClient, ids []string) error {
	for _, id := range ids {
		if _, err := client.DeleteMachine(ctx, &sri.DeleteMachineRequest{
			MachineId: id,
		}); err != nil {
			if status.Code(err) != codes.NotFound {
				return fmt.Errorf("error deleting machine %s: %w", id, err)
			}

			_, _ = fmt.Fprintf(streams.Out, "Machine %s not found\n", id)
		} else {
			_, _ = fmt.Fprintf(streams.Out, "Machine %s deleted\n", id)
		}
	}
	return nil
}
