// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package srictlmachine

import (
	goflag "flag"

	"github.com/spf13/cobra"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"spheric.cloud/spheric/srictl-machine/cmd/srictl-machine/srictlmachine/attach"
	"spheric.cloud/spheric/srictl-machine/cmd/srictl-machine/srictlmachine/common"
	"spheric.cloud/spheric/srictl-machine/cmd/srictl-machine/srictlmachine/create"
	"spheric.cloud/spheric/srictl-machine/cmd/srictl-machine/srictlmachine/delete"
	"spheric.cloud/spheric/srictl-machine/cmd/srictl-machine/srictlmachine/detach"
	"spheric.cloud/spheric/srictl-machine/cmd/srictl-machine/srictlmachine/exec"
	"spheric.cloud/spheric/srictl-machine/cmd/srictl-machine/srictlmachine/get"
	"spheric.cloud/spheric/srictl-machine/cmd/srictl-machine/srictlmachine/update"
	clicommon "spheric.cloud/spheric/srictl/cmd"
)

func Command(streams clicommon.Streams) *cobra.Command {
	var (
		zapOpts    zap.Options
		clientOpts common.Options
	)

	cmd := &cobra.Command{
		Use: "srictl-machine",
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
