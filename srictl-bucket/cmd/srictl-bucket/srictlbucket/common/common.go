// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package common

import (
	"fmt"
	"time"

	"github.com/spf13/pflag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	sri "spheric.cloud/spheric/sri/apis/bucket/v1alpha1"
	sriremotebucket "spheric.cloud/spheric/sri/remote/bucket"
	"spheric.cloud/spheric/srictl-bucket/renderers"
	srictlcmd "spheric.cloud/spheric/srictl/cmd"
	"spheric.cloud/spheric/srictl/renderer"
)

var Renderer = renderer.NewRegistry()

func init() {
	if err := renderers.AddToRegistry(Renderer); err != nil {
		panic(err)
	}
}

type ClientFactory interface {
	New() (sri.BucketRuntimeClient, func() error, error)
}

type ClientOptions struct {
	Address string
}

func (o *ClientOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.Address, "address", "", "Address to the sri server.")
}

func (o *ClientOptions) New() (sri.BucketRuntimeClient, func() error, error) {
	address, err := sriremotebucket.GetAddressWithTimeout(3*time.Second, o.Address)
	if err != nil {
		return nil, nil, err
	}

	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, fmt.Errorf("error dialing: %w", err)
	}

	return sri.NewBucketRuntimeClient(conn), conn.Close, nil
}

func NewOutputOptions() *srictlcmd.OutputOptions {
	return &srictlcmd.OutputOptions{
		Registry: Renderer,
	}
}

var (
	BucketAliases      = []string{"buckets"}
	BucketClassAliases = []string{"bucketchineclasses"}
)
