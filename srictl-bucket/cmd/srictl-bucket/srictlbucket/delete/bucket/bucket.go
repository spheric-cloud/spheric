// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package bucket

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	ctrl "sigs.k8s.io/controller-runtime"
	sri "spheric.cloud/spheric/sri/apis/bucket/v1alpha1"
	"spheric.cloud/spheric/srictl-bucket/cmd/srictl-bucket/srictlbucket/common"
	srictlcmd "spheric.cloud/spheric/srictl/cmd"
)

func Command(streams srictlcmd.Streams, clientFactory common.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "bucket id [ids...]",
		Aliases: common.BucketAliases,
		Args:    cobra.MinimumNArgs(1),
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

			ids := args

			return Run(cmd.Context(), streams, client, ids)
		},
	}

	return cmd
}

func Run(ctx context.Context, streams srictlcmd.Streams, client sri.BucketRuntimeClient, ids []string) error {
	for _, id := range ids {
		if _, err := client.DeleteBucket(ctx, &sri.DeleteBucketRequest{
			BucketId: id,
		}); err != nil {
			if status.Code(err) != codes.NotFound {
				return fmt.Errorf("error deleting bucket %s: %w", id, err)
			}

			_, _ = fmt.Fprintf(streams.Out, "Bucket %s not found\n", id)
		} else {
			_, _ = fmt.Fprintf(streams.Out, "Bucket %s deleted\n", id)
		}
	}
	return nil
}
