// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package common

import (
	"fmt"
	"time"

	"spheric.cloud/spheric/irictl/clientcmd"
	"spheric.cloud/spheric/irictl/tableconverters"

	iriremote "spheric.cloud/spheric/spherelet/iri/remote"

	"github.com/spf13/pflag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	iri "spheric.cloud/spheric/iri-api/apis/runtime/v1alpha1"
	"spheric.cloud/spheric/irictl/renderer"
	"spheric.cloud/spheric/irictl/tableconverter"
)

type Factory interface {
	Client() (iri.RuntimeServiceClient, func() error, error)
	Config() (*clientcmd.Config, error)
	Registry() (*renderer.Registry, error)
	OutputOptions() *OutputOptions
}

type Options struct {
	Address    string
	ConfigFile string
}

func NewOptions() *Options {
	return &Options{
		Address:    clientcmd.RecommendedInstanceRuntimeEndpoint,
		ConfigFile: "",
	}
}

func (o *Options) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.ConfigFile, clientcmd.RecommendedConfigPathFlag, o.ConfigFile, "Config file to use")
	fs.StringVar(&o.Address, "address", o.Address, "Address to the iri server.")
}

func (o *Options) Config() (*clientcmd.Config, error) {
	return clientcmd.GetConfig(o.ConfigFile)
}

func (o *Options) Registry() (*renderer.Registry, error) {
	registry := renderer.NewRegistry()
	if err := renderer.AddToRegistry(registry); err != nil {
		return nil, err
	}

	tableConvRegistry := tableconverter.NewRegistry()
	if err := tableconverters.AddToRegistry(tableConvRegistry); err != nil {
		return nil, err
	}

	if err := registry.Register("table", renderer.NewTable(tableConvRegistry)); err != nil {
		return nil, err
	}

	return registry, nil
}

func (o *Options) Client() (iri.RuntimeServiceClient, func() error, error) {
	address, err := iriremote.GetAddressWithTimeout(3*time.Second, o.Address)
	if err != nil {
		return nil, nil, err
	}

	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, fmt.Errorf("error dialing: %w", err)
	}

	return iri.NewRuntimeServiceClient(conn), conn.Close, nil
}

func (o *Options) OutputOptions() *OutputOptions {
	return &OutputOptions{
		factory: o,
	}
}

type OutputOptions struct {
	factory Factory
	Output  string
}

func (o *OutputOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVarP(&o.Output, "output", "o", o.Output, "Output format.")
}

func (o *OutputOptions) Renderer(ifEmpty string) (renderer.Renderer, error) {
	output := o.Output
	if output == "" {
		output = ifEmpty
	}

	r, err := o.factory.Registry()
	if err != nil {
		return nil, err
	}

	return r.Get(output)
}

func (o *OutputOptions) RendererOrNil() (renderer.Renderer, error) {
	output := o.Output
	if output == "" {
		return nil, nil
	}

	r, err := o.factory.Registry()
	if err != nil {
		return nil, err
	}

	return r.Get(output)
}

var (
	InstanceAliases         = []string{"instances", "inst", "insts"}
	DiskAliases             = []string{"disks", "dsk", "dsks"}
	NetworkInterfaceAliases = []string{"networkinterfaces", "nic", "nics"}
)
