// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"os"

	genericapiserver "k8s.io/apiserver/pkg/server"
	"k8s.io/component-base/cli"
	"spheric.cloud/spheric/internal/app/apiserver"
)

func main() {
	ctx := genericapiserver.SetupSignalContext()
	options := apiserver.NewSphericAPIServerOptions()
	cmd := apiserver.NewCommandStartSphericAPIServer(ctx, options)
	code := cli.Run(cmd)
	os.Exit(code)
}
