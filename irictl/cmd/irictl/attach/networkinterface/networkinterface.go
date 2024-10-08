// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package networkinterface

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	iri "spheric.cloud/spheric/iri-api/apis/runtime/v1alpha1"
	clicommon "spheric.cloud/spheric/irictl/cmd"
	"spheric.cloud/spheric/irictl/cmd/irictl/common"
	"spheric.cloud/spheric/irictl/decoder"
)

type Options struct {
	Filename   string
	InstanceID string
}

func (o *Options) AddFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&o.Filename, "filename", "f", o.Filename, "Path to a file to read.")
	cmd.Flags().StringVar(&o.InstanceID, "instance-id", "", "The instance ID to modify.")
	utilruntime.Must(cmd.MarkFlagRequired("instance-id"))
}

func Command(streams clicommon.Streams, clientFactory common.Factory) *cobra.Command {
	var (
		opts Options
	)

	cmd := &cobra.Command{
		Use:     "networkinterface",
		Aliases: common.NetworkInterfaceAliases,
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

			return Run(ctx, streams, client, opts)
		},
	}

	opts.AddFlags(cmd)

	return cmd
}

func Run(ctx context.Context, streams clicommon.Streams, client iri.RuntimeServiceClient, opts Options) error {
	data, err := clicommon.ReadFileOrReader(opts.Filename, os.Stdin)
	if err != nil {
		return err
	}

	networkInterface := &iri.NetworkInterface{}
	if err := decoder.Decode(data, networkInterface); err != nil {
		return err
	}

	if _, err := client.AttachNetworkInterface(ctx, &iri.AttachNetworkInterfaceRequest{NetworkInterface: networkInterface}); err != nil {
		return err
	}

	_, _ = fmt.Fprintf(streams.Out, "Attached networkinterface %s to instance %s\n", networkInterface.Name, opts.InstanceID)
	return nil
}
