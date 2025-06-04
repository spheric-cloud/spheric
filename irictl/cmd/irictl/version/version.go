// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package version

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	ctrl "sigs.k8s.io/controller-runtime"
	iri "spheric.cloud/spheric/iri-api/apis/runtime/v1alpha1"
	clicommon "spheric.cloud/spheric/irictl/cmd"
	"spheric.cloud/spheric/irictl/cmd/irictl/common"
)

func Command(streams clicommon.Streams, clientFactory common.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use: "version",
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

			return Run(ctx, streams, client)
		},
	}

	cmd.AddCommand()

	return cmd
}

func Run(ctx context.Context, streams clicommon.Streams, client iri.RuntimeServiceClient) error {
	res, err := client.Version(ctx, &iri.VersionRequest{})
	if err != nil {
		return fmt.Errorf("error getting version: %w", err)
	}

	_, _ = fmt.Fprintln(streams.Out, "Runtime Name", res.RuntimeName)
	_, _ = fmt.Fprintln(streams.Out, "Runtime Version", res.RuntimeVersion)
	return nil
}
