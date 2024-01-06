// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package bucketclass

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	ctrl "sigs.k8s.io/controller-runtime"
	sri "spheric.cloud/spheric/sri/apis/bucket/v1alpha1"
	"spheric.cloud/spheric/srictl-bucket/cmd/srictl-bucket/srictlbucket/common"
	srictlcmd "spheric.cloud/spheric/srictl/cmd"
	"spheric.cloud/spheric/srictl/renderer"
)

func Command(streams srictlcmd.Streams, clientFactory common.ClientFactory) *cobra.Command {
	var (
		outputOpts = common.NewOutputOptions()
	)

	cmd := &cobra.Command{
		Use:     "bucketclass",
		Aliases: common.BucketClassAliases,
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

			return Run(cmd.Context(), streams, client, render)
		},
	}

	outputOpts.AddFlags(cmd.Flags())

	return cmd
}

func Run(ctx context.Context, streams srictlcmd.Streams, client sri.BucketRuntimeClient, render renderer.Renderer) error {
	res, err := client.ListBucketClasses(ctx, &sri.ListBucketClassesRequest{})
	if err != nil {
		return fmt.Errorf("error listing bucket classes: %w", err)
	}

	return render.Render(res.BucketClasses, streams.Out)
}
