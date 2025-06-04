// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0

package vee

import (
	"context"
	"fmt"
	"net"
	"os/user"
	"path/filepath"

	"github.com/go-logr/logr"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"spheric.cloud/spheric/actuo/codec"
	"spheric.cloud/spheric/actuo/etcd/embed"
	"spheric.cloud/spheric/actuo/run"
	"spheric.cloud/spheric/actuo/storage/etcd"
	storagestore "spheric.cloud/spheric/actuo/storage/store"
	iri "spheric.cloud/spheric/iri-api/apis/runtime/v1alpha1"
	utilgrpc "spheric.cloud/spheric/utils/grpc"
	utilos "spheric.cloud/spheric/utils/os"
	"spheric.cloud/spheric/vee/api"
	"spheric.cloud/spheric/vee/iriserver"
	"spheric.cloud/spheric/vee/server"
)

var (
	defaultDir string
)

func init() {
	u, err := user.Current()
	utilruntime.Must(err)
	defaultDir = filepath.Join(u.HomeDir, ".vee")
}

type Options struct {
	APISocket string
	Dir       string
}

func NewOptions() *Options {
	return &Options{
		APISocket: filepath.Join("/var", "run", "vee", "vee.sock"),
		Dir:       defaultDir,
	}
}

func (o *Options) AddFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&o.APISocket, "api-socket", o.APISocket, "Where to create the API socket.")
	cmd.Flags().StringVar(&o.Dir, "dir", o.Dir, "Directory to store data in.")
}

func Command() *cobra.Command {
	var (
		zapOpts = zap.Options{Development: true}
		opts    = NewOptions()
	)

	cmd := &cobra.Command{
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			logger := zap.New(zap.UseFlagOptions(&zapOpts))
			ctrl.SetLogger(logger)
			cmd.SetContext(ctrl.LoggerInto(cmd.Context(), ctrl.Log))
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			return Run(ctx, *opts)
		},
	}

	opts.AddFlags(cmd)

	return cmd
}

func startGRPCServer(ctx context.Context, setupLog logr.Logger, dir, apiSocket string) error {
	srv, err := iriserver.New(dir)
	if err != nil {
		return fmt.Errorf("error creating server: %w", err)
	}

	grpcSrv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			utilgrpc.InjectLogger(ctrl.Log.WithName("server")),
			utilgrpc.LogRequest,
		),
	)

	if err := utilos.EnsureSocketGone(apiSocket); err != nil {
		return err
	}
	l, err := net.Listen("unix", apiSocket)
	if err != nil {
		return fmt.Errorf("unable to create listener: %w", err)
	}
	defer func() { _ = l.Close() }()

	iri.RegisterRuntimeServiceServer(grpcSrv, srv)

	done := make(chan struct{})
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		defer close(done)
		defer cancel()

		setupLog.Info("Starting grpc server", "Address", l.Addr().String())
		_ = grpcSrv.Serve(l)
		setupLog.Info("Stopping grpc server")
	}()

	<-ctx.Done()
	grpcSrv.GracefulStop()
	<-done
	return nil
}

func Run(ctx context.Context, opts Options) error {
	setupLog := ctrl.Log.WithName("setup")

	e, err := embed.New(embed.WithLogger(ctrl.Log.WithName("etcd")))
	if err != nil {
		return err
	}
	if err := e.Start(); err != nil {
		return err
	}
	defer func() { _ = e.Stop() }()

	c, err := e.NewClient()
	if err != nil {
		return err
	}
	defer func() { _ = c.Close() }()

	store := etcd.NewSimple[*api.Instance](
		c,
		codec.JSON[*api.Instance](),
		storagestore.DefaultFactory[*api.Instance](),
		storagestore.DefaultMetaVersioner[*api.Instance](),
	)
	veeSrv, err := server.New("vee.sock", store)
	if err != nil {
		return err
	}

	g := run.NewGroup(ctx)

	g.Start(veeSrv.ListenAndServe, run.OnErrorStop)
	g.Start(func(ctx context.Context) error {
		return startGRPCServer(ctx, setupLog, opts.Dir, opts.APISocket)
	}, run.OnErrorStop)

	return g.Wait()
}
