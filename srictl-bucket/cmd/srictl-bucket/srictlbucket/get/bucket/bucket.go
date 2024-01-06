// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package bucket

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	ctrl "sigs.k8s.io/controller-runtime"
	sri "spheric.cloud/spheric/sri/apis/bucket/v1alpha1"
	"spheric.cloud/spheric/srictl-bucket/cmd/srictl-bucket/srictlbucket/common"
	srictlcmd "spheric.cloud/spheric/srictl/cmd"
	"spheric.cloud/spheric/srictl/renderer"
)

type Options struct {
}

func (o *Options) AddFlags(fs *pflag.FlagSet) {
}

func Command(streams srictlcmd.Streams, clientFactory common.ClientFactory) *cobra.Command {
	var (
		opts       Options
		outputOpts = common.NewOutputOptions()
	)

	cmd := &cobra.Command{
		Use:     "bucket",
		Aliases: common.BucketAliases,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			log := ctrl.LoggerFrom(ctx)

			client, cleanup, err := clientFactory.New()
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

			return Run(cmd.Context(), streams, client, render, opts)
		},
	}

	outputOpts.AddFlags(cmd.Flags())
	opts.AddFlags(cmd.Flags())

	return cmd
}

func Run(ctx context.Context, streams srictlcmd.Streams, client sri.BucketRuntimeClient, render renderer.Renderer, opts Options) error {
	res, err := client.ListBuckets(ctx, &sri.ListBucketsRequest{})
	if err != nil {
		return fmt.Errorf("error listing buckets: %w", err)
	}

	return render.Render(res.Buckets, streams.Out)
}
