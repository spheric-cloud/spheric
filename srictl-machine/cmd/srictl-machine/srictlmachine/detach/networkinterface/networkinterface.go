// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package networkinterface

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	sri "spheric.cloud/spheric/sri/apis/machine/v1alpha1"
	"spheric.cloud/spheric/srictl-machine/cmd/srictl-machine/srictlmachine/common"
	clicommon "spheric.cloud/spheric/srictl/cmd"
)

type Options struct {
	MachineID string
}

func (o *Options) AddFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&o.MachineID, "machine-id", "", "The machine ID to modify.")
	utilruntime.Must(cmd.MarkFlagRequired("machine-id"))
}

func Command(streams clicommon.Streams, clientFactory common.Factory) *cobra.Command {
	var (
		opts Options
	)

	cmd := &cobra.Command{
		Use:     "networkinterface name [names...]",
		Aliases: common.NetworkInterfaceAliases,
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

			names := args
			return Run(ctx, streams, client, names, opts)
		},
	}

	opts.AddFlags(cmd)

	return cmd
}

func Run(ctx context.Context, streams clicommon.Streams, client sri.MachineRuntimeClient, names []string, opts Options) error {
	for _, name := range names {
		if _, err := client.DetachNetworkInterface(ctx, &sri.DetachNetworkInterfaceRequest{
			MachineId: opts.MachineID,
			Name:      name,
		}); err != nil {
			if status.Code(err) != codes.NotFound {
				return fmt.Errorf("error detaching network interface %s from machine %s: %w", name, opts.MachineID, err)
			}
			_, _ = fmt.Fprintf(streams.Out, "Network interface %s in machine %s not found\n", name, opts.MachineID)
		} else {
			_, _ = fmt.Fprintf(streams.Out, "Detached network interface %s from machine %s\n", name, opts.MachineID)
		}
	}
	return nil
}
