// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package bucket

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	ctrl "sigs.k8s.io/controller-runtime"
	sri "spheric.cloud/spheric/sri/apis/bucket/v1alpha1"
	"spheric.cloud/spheric/srictl-bucket/cmd/srictl-bucket/srictlbucket/common"
	srictlcmd "spheric.cloud/spheric/srictl/cmd"
	"spheric.cloud/spheric/srictl/decoder"
	"spheric.cloud/spheric/srictl/renderer"
)

type Options struct {
	Filename string
}

func (o *Options) AddFlags(fs *pflag.FlagSet) {
	fs.StringVarP(&o.Filename, "filename", "f", o.Filename, "Path to a file to read.")
}

func Command(streams srictlcmd.Streams, clientFactory common.ClientFactory) *cobra.Command {
	var (
		outputOpts = common.NewOutputOptions()
		opts       Options
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

			r, err := outputOpts.RendererOrNil()
			if err != nil {
				return err
			}

			return Run(ctx, streams, client, r, opts)
		},
	}

	outputOpts.AddFlags(cmd.Flags())
	opts.AddFlags(cmd.Flags())

	return cmd
}

func Run(ctx context.Context, streams srictlcmd.Streams, client sri.BucketRuntimeClient, r renderer.Renderer, opts Options) error {
	data, err := srictlcmd.ReadFileOrReader(opts.Filename, os.Stdin)
	if err != nil {
		return err
	}

	bucket := &sri.Bucket{}
	if err := decoder.Decode(data, bucket); err != nil {
		return err
	}

	res, err := client.CreateBucket(ctx, &sri.CreateBucketRequest{Bucket: bucket})
	if err != nil {
		return err
	}

	if r != nil {
		return r.Render(res.Bucket, streams.Out)
	}

	_, _ = fmt.Fprintf(streams.Out, "Created bucket %s\n", res.Bucket.Metadata.Id)
	return nil
}
