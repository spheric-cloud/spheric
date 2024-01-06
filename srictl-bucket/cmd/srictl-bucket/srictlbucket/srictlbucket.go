// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package srictlbucket

import (
	goflag "flag"

	"github.com/spf13/cobra"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"spheric.cloud/spheric/srictl-bucket/cmd/srictl-bucket/srictlbucket/common"
	"spheric.cloud/spheric/srictl-bucket/cmd/srictl-bucket/srictlbucket/create"
	delete2 "spheric.cloud/spheric/srictl-bucket/cmd/srictl-bucket/srictlbucket/delete"
	"spheric.cloud/spheric/srictl-bucket/cmd/srictl-bucket/srictlbucket/get"
	srictlcmd "spheric.cloud/spheric/srictl/cmd"
)

func Command(streams srictlcmd.Streams) *cobra.Command {
	var (
		zapOpts    zap.Options
		clientOpts common.ClientOptions
	)

	cmd := &cobra.Command{
		Use: "srictl-bucket",
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
		delete2.Command(streams, &clientOpts),
		create.Command(streams, &clientOpts),
	)

	return cmd
}
