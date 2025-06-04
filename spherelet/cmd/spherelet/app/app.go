// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package app

import (
	"context"
	goflag "flag"
	"fmt"
	"net"
	"strconv"
	"time"

	"spheric.cloud/spheric/spherelet/event/instanceevent"
	"spheric.cloud/spheric/spherelet/event/runtimeevent"

	"spheric.cloud/spheric/spherelet/iri/remote"

	"github.com/ironcore-dev/controller-utils/configutils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"
	corev1alpha1 "spheric.cloud/spheric/api/core/v1alpha1"
	coreclient "spheric.cloud/spheric/internal/client/core"
	iri "spheric.cloud/spheric/iri-api/apis/runtime/v1alpha1"
	"spheric.cloud/spheric/spherelet/addresses"
	sphereletclient "spheric.cloud/spheric/spherelet/client"
	sphereletclientconfig "spheric.cloud/spheric/spherelet/client/config"
	"spheric.cloud/spheric/spherelet/controllers"
	"spheric.cloud/spheric/spherelet/server"
	"spheric.cloud/spheric/utils/client/config"
)

var (
	scheme = runtime.NewScheme()
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(corev1alpha1.AddToScheme(scheme))
}

type Options struct {
	GetConfigOptions         config.GetConfigOptions
	MetricsAddr              string
	EnableLeaderElection     bool
	LeaderElectionNamespace  string
	LeaderElectionKubeconfig string
	ProbeAddr                string

	FleetName                             string
	InstanceDownwardAPILabels             map[string]string
	InstanceDownwardAPIAnnotations        map[string]string
	ProviderID                            string
	InstanceRuntimeEndpoint               string
	InstanceRuntimeSocketDiscoveryTimeout time.Duration
	DialTimeout                           time.Duration
	InstanceTypeMapperSyncTimeout         time.Duration

	ServerFlags server.Flags

	AddressesOptions addresses.GetOptions

	WatchFilterValue string
}

func (o *Options) AddFlags(fs *pflag.FlagSet) {
	o.GetConfigOptions.BindFlags(fs)
	fs.StringVar(&o.MetricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	fs.StringVar(&o.ProbeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	fs.BoolVar(&o.EnableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	fs.StringVar(&o.LeaderElectionNamespace, "leader-election-namespace", "", "Namespace to do leader election in.")
	fs.StringVar(&o.LeaderElectionKubeconfig, "leader-election-kubeconfig", "", "Path pointing to a kubeconfig to use for leader election.")

	fs.StringVar(&o.FleetName, "instance-pool-name", o.FleetName, "Name of the instance pool to announce / watch")
	fs.StringToStringVar(&o.InstanceDownwardAPILabels, "instance-downward-api-label", o.InstanceDownwardAPILabels, "Downward-API labels to set on the iri instance.")
	fs.StringToStringVar(&o.InstanceDownwardAPIAnnotations, "instance-downward-api-annotation", o.InstanceDownwardAPIAnnotations, "Downward-API annotations to set on the iri instance.")
	fs.StringVar(&o.ProviderID, "provider-id", "", "Provider id to announce on the instance pool.")
	fs.StringVar(&o.InstanceRuntimeEndpoint, "instance-runtime-endpoint", o.InstanceRuntimeEndpoint, "Endpoint of the remote instance runtime service.")
	fs.DurationVar(&o.InstanceRuntimeSocketDiscoveryTimeout, "instance-runtime-socket-discovery-timeout", 20*time.Second, "Timeout for discovering the instance runtime socket.")
	fs.DurationVar(&o.DialTimeout, "dial-timeout", 1*time.Second, "Timeout for dialing to the instance runtime endpoint.")

	o.ServerFlags.BindFlags(fs)

	o.AddressesOptions.BindFlags(fs)

	fs.StringVar(&o.WatchFilterValue, "watch-filter", "", "Value to filter for while watching.")
}

func (o *Options) MarkFlagsRequired(cmd *cobra.Command) {
	_ = cmd.MarkFlagRequired("instance-pool-name")
	_ = cmd.MarkFlagRequired("provider-id")
}

func NewOptions() *Options {
	return &Options{
		ServerFlags: *server.NewServerFlags(),
	}
}

func Command() *cobra.Command {
	var (
		zapOpts = zap.Options{Development: true}
		opts    = NewOptions()
	)

	cmd := &cobra.Command{
		Use: "spherelet",
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

	goFlags := goflag.NewFlagSet("", 0)
	zapOpts.BindFlags(goFlags)
	cmd.PersistentFlags().AddGoFlagSet(goFlags)

	opts.AddFlags(cmd.Flags())
	opts.MarkFlagsRequired(cmd)

	return cmd
}

func getPort(address string) (int32, error) {
	_, portString, err := net.SplitHostPort(address)
	if err != nil {
		return 0, fmt.Errorf("error splitting serving address into host / port: %w", err)
	}

	portInt64, err := strconv.ParseInt(portString, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("error parsing port %q: %w", portString, err)
	}

	if portInt64 == 0 {
		return 0, fmt.Errorf("cannot specify dynamic port")
	}
	return int32(portInt64), nil
}

func Run(ctx context.Context, opts Options) error {
	logger := ctrl.LoggerFrom(ctx)
	setupLog := ctrl.Log.WithName("setup")

	port, err := getPort(opts.ServerFlags.Serving.Address)
	if err != nil {
		return fmt.Errorf("error getting port from address: %w", err)
	}

	getter, err := sphereletclientconfig.NewGetter(opts.FleetName)
	if err != nil {
		return fmt.Errorf("error creating new getter: %w", err)
	}

	endpoint, err := remote.GetAddressWithTimeout(opts.InstanceRuntimeSocketDiscoveryTimeout, opts.InstanceRuntimeEndpoint)
	if err != nil {
		return fmt.Errorf("error detecting instance runtime endpoint: %w", err)
	}

	fleetAddresses, err := addresses.Get(&opts.AddressesOptions)
	if err != nil {
		return fmt.Errorf("error getting instance pool endpoints: %w", err)
	}

	setupLog.V(1).Info("Discovered addresses to report", "FleetAddresses", fleetAddresses)

	instanceRuntime, err := remote.NewRemoteRuntime(endpoint)
	if err != nil {
		return fmt.Errorf("error creating remote instance runtime: %w", err)
	}

	cfg, configCtrl, err := getter.GetConfig(ctx, &opts.GetConfigOptions)
	if err != nil {
		return fmt.Errorf("error getting config: %w", err)
	}

	leaderElectionCfg, err := configutils.GetConfig(
		configutils.Kubeconfig(opts.LeaderElectionKubeconfig),
	)
	if err != nil {
		return fmt.Errorf("error creating leader election kubeconfig: %w", err)
	}

	mgr, err := ctrl.NewManager(cfg, ctrl.Options{
		Logger:                  logger,
		Scheme:                  scheme,
		Metrics:                 metricsserver.Options{BindAddress: opts.MetricsAddr},
		HealthProbeBindAddress:  opts.ProbeAddr,
		LeaderElection:          opts.EnableLeaderElection,
		LeaderElectionID:        "bfafcebe.spheric.cloud",
		LeaderElectionNamespace: opts.LeaderElectionNamespace,
		LeaderElectionConfig:    leaderElectionCfg,
		Cache:                   cache.Options{ByObject: map[client.Object]cache.ByObject{}},
		NewCache: func(config *rest.Config, cacheOpts cache.Options) (cache.Cache, error) {
			cacheOpts.ByObject[&corev1alpha1.Instance{}] = cache.ByObject{
				Field: fields.OneTermEqualSelector(
					corev1alpha1.InstanceFleetRefNameField,
					opts.FleetName,
				),
			}
			return cache.New(config, cacheOpts)
		},
	})
	if err != nil {
		return fmt.Errorf("error creating manager: %w", err)
	}
	if err := config.SetupControllerWithManager(mgr, configCtrl); err != nil {
		return err
	}

	version, err := instanceRuntime.Version(ctx, &iri.VersionRequest{})
	if err != nil {
		return fmt.Errorf("error getting instance runtime version: %w", err)
	}

	srvOpts := opts.ServerFlags.ServerOptions(
		opts.FleetName,
		instanceRuntime,
		logger.WithName("server"),
	)
	srv, err := server.New(cfg, srvOpts)
	if err != nil {
		return fmt.Errorf("error creating spherelet server: %w", err)
	}

	if err := mgr.Add(srv); err != nil {
		return fmt.Errorf("error adding spherelet server to manager: %w", err)
	}

	instanceEvents := instanceevent.NewGenerator(func(ctx context.Context) ([]*iri.Instance, error) {
		res, err := instanceRuntime.ListInstances(ctx, &iri.ListInstancesRequest{})
		if err != nil {
			return nil, err
		}
		return res.Instances, nil
	}, instanceevent.GeneratorOptions{})
	if err := mgr.Add(instanceEvents); err != nil {
		return fmt.Errorf("error adding instance event generator: %w", err)
	}
	if err := mgr.AddHealthzCheck("instance-events", instanceEvents.Check); err != nil {
		return fmt.Errorf("error adding instance event generator healthz check: %w", err)
	}

	runtimeEvents := runtimeevent.NewGenerator(func(ctx context.Context) (*iri.RuntimeResources, error) {
		res, err := instanceRuntime.Status(ctx, &iri.StatusRequest{})
		if err != nil {
			return nil, err
		}
		return res.Allocatable, nil
	}, runtimeevent.GeneratorOptions{})
	if err := mgr.Add(runtimeEvents); err != nil {
		return fmt.Errorf("error adding runtime event generator: %w", err)
	}
	if err := mgr.AddHealthzCheck("runtime-events", runtimeEvents.Check); err != nil {
		return fmt.Errorf("error adding runtime event generator healthz check: %w", err)
	}

	indexer := mgr.GetFieldIndexer()
	if err := sphereletclient.SetupInstanceSpecSecretNamesField(ctx, indexer, opts.FleetName); err != nil {
		return fmt.Errorf("error setting up %s indexer with manager: %w", sphereletclient.InstanceSpecSecretNamesField, err)
	}
	if err := sphereletclient.SetupInstanceSpecDiskNamesField(ctx, indexer, opts.FleetName); err != nil {
		return fmt.Errorf("error setting up %s indexer with manager: %w", sphereletclient.InstanceSpecDiskNamesField, err)
	}

	if err := coreclient.SetupInstanceSpecFleetRefNameFieldIndexer(ctx, indexer); err != nil {
		return fmt.Errorf("error setting up %s indexer with manager: %w", coreclient.InstanceSpecFleetRefNameField, err)
	}

	onInitialized := func(ctx context.Context) error {
		if err := (&controllers.InstanceReconciler{
			EventRecorder:          mgr.GetEventRecorderFor("instances"),
			Client:                 mgr.GetClient(),
			InstanceRuntime:        instanceRuntime,
			InstanceRuntimeName:    version.RuntimeName,
			InstanceRuntimeVersion: version.RuntimeVersion,
			FleetName:              opts.FleetName,
			DownwardAPILabels:      opts.InstanceDownwardAPILabels,
			DownwardAPIAnnotations: opts.InstanceDownwardAPIAnnotations,
			WatchFilterValue:       opts.WatchFilterValue,
		}).SetupWithManager(mgr); err != nil {
			return fmt.Errorf("error setting up instance reconciler with manager: %w", err)
		}

		if err := (&controllers.InstanceAnnotatorReconciler{
			Client:         mgr.GetClient(),
			InstanceEvents: instanceEvents,
		}).SetupWithManager(mgr); err != nil {
			return fmt.Errorf("error setting up instance annotator reconciler with manager: %w", err)
		}

		if err := (&controllers.FleetReconciler{
			Client:          mgr.GetClient(),
			FleetName:       opts.FleetName,
			Addresses:       fleetAddresses,
			Port:            port,
			InstanceRuntime: instanceRuntime,
		}).SetupWithManager(mgr); err != nil {
			return fmt.Errorf("error setting up instance pool reconciler with manager: %w", err)
		}

		if err := (&controllers.FleetAnnotatorReconciler{
			Client:        mgr.GetClient(),
			FleetName:     opts.FleetName,
			RuntimeEvents: runtimeEvents,
		}).SetupWithManager(mgr); err != nil {
			return fmt.Errorf("error setting up instance pool annotator reconciler with manager: %w", err)
		}

		return nil
	}

	if err := (&controllers.FleetInit{
		Client:        mgr.GetClient(),
		FleetName:     opts.FleetName,
		ProviderID:    opts.ProviderID,
		OnInitialized: onInitialized,
	}).SetupWithManager(mgr); err != nil {
		return fmt.Errorf("error setting up instance pool init with manager: %w", err)
	}

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		return fmt.Errorf("error adding healthz check: %w", err)
	}

	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		return fmt.Errorf("error adding readyz check: %w", err)
	}

	setupLog.Info("starting manager")
	if err := mgr.Start(ctx); err != nil {
		return fmt.Errorf("error running manager: %w", err)
	}
	return nil
}
