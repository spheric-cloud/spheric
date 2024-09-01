// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package irictl

import (
	goflag "flag"

	"github.com/spf13/cobra"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	clicommon "spheric.cloud/spheric/irictl/cmd"
	"spheric.cloud/spheric/irictl/cmd/irictl/attach"
	"spheric.cloud/spheric/irictl/cmd/irictl/common"
	"spheric.cloud/spheric/irictl/cmd/irictl/create"
	"spheric.cloud/spheric/irictl/cmd/irictl/delete"
	"spheric.cloud/spheric/irictl/cmd/irictl/detach"
	"spheric.cloud/spheric/irictl/cmd/irictl/exec"
	"spheric.cloud/spheric/irictl/cmd/irictl/get"
	"spheric.cloud/spheric/irictl/cmd/irictl/update"
)

func Command(streams clicommon.Streams) *cobra.Command {
	var (
		zapOpts    zap.Options
		clientOpts common.Options
	)

	cmd := &cobra.Command{
		Use: "irictl-instance",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			logger := zap.New(zap.UseFlagOptions(&zapOpts))
			ctrl.SetLogger(logger)
			cmd.SetContext(ctrl.LoggerInto(cmd.Context(), ctrl.Log))
		},
	}

	goFlags := goflag.NewFlagSet("", 0)
	zapOpts.BindFlags(goFlags)

	cmd.PersistentFlags().AddGoFlagSet(goFlags)
	clientOpts.AddFlags(cmd.PersistentFlags())

	cmd.AddCommand(
		get.Command(streams, &clientOpts),
		create.Command(streams, &clientOpts),
		delete.Command(streams, &clientOpts),
		update.Command(streams, &clientOpts),
		exec.Command(streams, &clientOpts),
		attach.Command(streams, &clientOpts),
		detach.Command(streams, &clientOpts),
	)

	return cmd
}
