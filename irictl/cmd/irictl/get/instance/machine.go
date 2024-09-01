// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package instance

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	ctrl "sigs.k8s.io/controller-runtime"
	sri "spheric.cloud/spheric/iri-api/apis/runtime/v1alpha1"
	clicommon "spheric.cloud/spheric/irictl/cmd"
	"spheric.cloud/spheric/irictl/cmd/irictl/common"
	"spheric.cloud/spheric/irictl/renderer"
)

type Options struct {
	Labels map[string]string
}

func (o *Options) AddFlags(fs *pflag.FlagSet) {
	fs.StringToStringVarP(&o.Labels, "labels", "l", o.Labels, "Labels to filter the instances by.")
}

func Command(streams clicommon.Streams, clientFactory common.Factory) *cobra.Command {
	var (
		opts       Options
		outputOpts = clientFactory.OutputOptions()
	)

	cmd := &cobra.Command{
		Use:     "instance name",
		Args:    cobra.MaximumNArgs(1),
		Aliases: common.InstanceAliases,
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

			var name string
			if len(args) > 0 {
				name = args[0]
			}

			return Run(cmd.Context(), streams, client, render, name, opts)
		},
	}

	outputOpts.AddFlags(cmd.Flags())
	opts.AddFlags(cmd.Flags())

	return cmd
}

func Run(
	ctx context.Context,
	streams clicommon.Streams,
	client sri.RuntimeServiceClient,
	render renderer.Renderer,
	name string,
	opts Options,
) error {
	var filter *sri.InstanceFilter
	if name != "" || opts.Labels != nil {
		filter = &sri.InstanceFilter{
			Id:            name,
			LabelSelector: opts.Labels,
		}
	}

	res, err := client.ListInstances(ctx, &sri.ListInstancesRequest{Filter: filter})
	if err != nil {
		return fmt.Errorf("error listing instances: %w", err)
	}

	return render.Render(res.Instances, streams.Out)
}
