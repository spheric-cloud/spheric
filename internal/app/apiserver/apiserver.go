// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package apiserver

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apiserver/pkg/admission"
	"k8s.io/apiserver/pkg/endpoints/openapi"
	"k8s.io/apiserver/pkg/features"
	genericapiserver "k8s.io/apiserver/pkg/server"
	genericoptions "k8s.io/apiserver/pkg/server/options"
	serverstorage "k8s.io/apiserver/pkg/server/storage"
	utilfeature "k8s.io/apiserver/pkg/util/feature"
	netutils "k8s.io/utils/net"
	computev1alpha1 "spheric.cloud/spheric/api/compute/v1alpha1"
	corev1alpha1 "spheric.cloud/spheric/api/core/v1alpha1"
	ipamv1alpha1 "spheric.cloud/spheric/api/ipam/v1alpha1"
	networkingv1alpha1 "spheric.cloud/spheric/api/networking/v1alpha1"
	storagev1alpha1 "spheric.cloud/spheric/api/storage/v1alpha1"
	"spheric.cloud/spheric/client-go/informers"
	sphericopenapi "spheric.cloud/spheric/client-go/openapi"
	clientset "spheric.cloud/spheric/client-go/spheric"
	sphericinitializer "spheric.cloud/spheric/internal/admission/initializer"
	"spheric.cloud/spheric/internal/admission/plugin/machinevolumedevices"
	"spheric.cloud/spheric/internal/admission/plugin/resourcequota"
	"spheric.cloud/spheric/internal/admission/plugin/volumeresizepolicy"
	"spheric.cloud/spheric/internal/api"
	"spheric.cloud/spheric/internal/apis/compute"
	"spheric.cloud/spheric/internal/apiserver"
	"spheric.cloud/spheric/internal/machinepoollet/client"
	"spheric.cloud/spheric/internal/quota/evaluator/spheric"
	apiequality "spheric.cloud/spheric/utils/equality"
	"spheric.cloud/spheric/utils/quota"
)

const defaultEtcdPathPrefix = "/registry/spheric.cloud"

func init() {
	utilruntime.Must(apiequality.AddFuncs(equality.Semantic))
}

func NewResourceConfig() *serverstorage.ResourceConfig {
	cfg := serverstorage.NewResourceConfig()
	cfg.EnableVersions(
		computev1alpha1.SchemeGroupVersion,
		corev1alpha1.SchemeGroupVersion,
		storagev1alpha1.SchemeGroupVersion,
		networkingv1alpha1.SchemeGroupVersion,
		ipamv1alpha1.SchemeGroupVersion,
	)
	return cfg
}

type SphericAPIServerOptions struct {
	RecommendedOptions   *genericoptions.RecommendedOptions
	MachinePoolletConfig client.MachinePoolletClientConfig

	SharedInformerFactory informers.SharedInformerFactory
}

func (o *SphericAPIServerOptions) AddFlags(fs *pflag.FlagSet) {
	o.RecommendedOptions.AddFlags(fs)

	// machinepoollet related flags:
	fs.StringSliceVar(&o.MachinePoolletConfig.PreferredAddressTypes, "machinepoollet-preferred-address-types", o.MachinePoolletConfig.PreferredAddressTypes,
		"List of the preferred MachinePoolAddressTypes to use for machinepoollet connections.")

	fs.DurationVar(&o.MachinePoolletConfig.HTTPTimeout, "machinepoollet-timeout", o.MachinePoolletConfig.HTTPTimeout,
		"Timeout for machinepoollet operations.")

	fs.StringVar(&o.MachinePoolletConfig.CertFile, "machinepoollet-client-certificate", o.MachinePoolletConfig.CertFile,
		"Path to a client cert file for TLS.")

	fs.StringVar(&o.MachinePoolletConfig.KeyFile, "machinepoollet-client-key", o.MachinePoolletConfig.KeyFile,
		"Path to a client key file for TLS.")

	fs.StringVar(&o.MachinePoolletConfig.CAFile, "machinepoollet-certificate-authority", o.MachinePoolletConfig.CAFile,
		"Path to a cert file for the certificate authority.")
}

func NewSphericAPIServerOptions() *SphericAPIServerOptions {
	o := &SphericAPIServerOptions{
		RecommendedOptions: genericoptions.NewRecommendedOptions(
			defaultEtcdPathPrefix,
			api.Codecs.LegacyCodec(
				computev1alpha1.SchemeGroupVersion,
				corev1alpha1.SchemeGroupVersion,
				storagev1alpha1.SchemeGroupVersion,
				networkingv1alpha1.SchemeGroupVersion,
				ipamv1alpha1.SchemeGroupVersion,
			),
		),
		MachinePoolletConfig: client.MachinePoolletClientConfig{
			Port:         12319,
			ReadOnlyPort: 12320,
			PreferredAddressTypes: []string{
				string(compute.MachinePoolHostName),

				// internal, preferring DNS if reported
				string(compute.MachinePoolInternalDNS),
				string(compute.MachinePoolInternalIP),

				// external, preferring DNS if reported
				string(compute.MachinePoolExternalDNS),
				string(compute.MachinePoolExternalIP),
			},
			HTTPTimeout: time.Duration(5) * time.Second,
		},
	}
	o.RecommendedOptions.Etcd.StorageConfig.EncodeVersioner = runtime.NewMultiGroupVersioner(
		computev1alpha1.SchemeGroupVersion,
		schema.GroupKind{Group: computev1alpha1.SchemeGroupVersion.Group},
		schema.GroupKind{Group: corev1alpha1.SchemeGroupVersion.Group},
		schema.GroupKind{Group: storagev1alpha1.SchemeGroupVersion.Group},
		schema.GroupKind{Group: networkingv1alpha1.SchemeGroupVersion.Group},
		schema.GroupKind{Group: ipamv1alpha1.SchemeGroupVersion.Group},
	)
	return o
}

func NewCommandStartSphericAPIServer(ctx context.Context, defaults *SphericAPIServerOptions) *cobra.Command {
	o := *defaults
	cmd := &cobra.Command{
		Short: "Launch an sphericAPI server",
		Long:  "Launch an sphericAPI server",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.Complete(); err != nil {
				return err
			}
			if err := o.Validate(args); err != nil {
				return err
			}
			if err := o.Run(ctx); err != nil {
				return err
			}
			return nil
		},
	}

	o.AddFlags(cmd.Flags())
	utilfeature.DefaultMutableFeatureGate.AddFlag(cmd.Flags())

	return cmd
}

func (o *SphericAPIServerOptions) Validate(args []string) error {
	var errors []error
	errors = append(errors, o.RecommendedOptions.Validate()...)
	return utilerrors.NewAggregate(errors)
}

func (o *SphericAPIServerOptions) Complete() error {
	machinevolumedevices.Register(o.RecommendedOptions.Admission.Plugins)
	resourcequota.Register(o.RecommendedOptions.Admission.Plugins)
	volumeresizepolicy.Register(o.RecommendedOptions.Admission.Plugins)

	o.RecommendedOptions.Admission.RecommendedPluginOrder = append(
		o.RecommendedOptions.Admission.RecommendedPluginOrder,
		machinevolumedevices.PluginName,
		resourcequota.PluginName,
		volumeresizepolicy.PluginName,
	)

	return nil
}

func (o *SphericAPIServerOptions) Config() (*apiserver.Config, error) {
	if err := o.RecommendedOptions.SecureServing.MaybeDefaultWithSelfSignedCerts("localhost", nil, []net.IP{netutils.ParseIPSloppy("127.0.0.1")}); err != nil {
		return nil, fmt.Errorf("error creating self-signed certificates: %w", err)
	}

	o.RecommendedOptions.Etcd.StorageConfig.Paging = utilfeature.DefaultFeatureGate.Enabled(features.APIListChunking)

	o.RecommendedOptions.ExtraAdmissionInitializers = func(c *genericapiserver.RecommendedConfig) ([]admission.PluginInitializer, error) {
		sphericClient, err := clientset.NewForConfig(c.LoopbackClientConfig)
		if err != nil {
			return nil, err
		}

		informerFactory := informers.NewSharedInformerFactory(sphericClient, c.LoopbackClientConfig.Timeout)
		o.SharedInformerFactory = informerFactory

		quotaRegistry := quota.NewRegistry(api.Scheme)
		if err := quota.AddAllToRegistry(quotaRegistry, spheric.NewEvaluatorsForAdmission(sphericClient, informerFactory)); err != nil {
			return nil, fmt.Errorf("error initializing quota registry: %w", err)
		}

		genericInitializer := sphericinitializer.New(sphericClient, informerFactory, quotaRegistry)

		return []admission.PluginInitializer{
			genericInitializer,
		}, nil
	}

	serverConfig := genericapiserver.NewRecommendedConfig(api.Codecs)

	serverConfig.OpenAPIConfig = genericapiserver.DefaultOpenAPIConfig(sphericopenapi.GetOpenAPIDefinitions, openapi.NewDefinitionNamer(api.Scheme))
	serverConfig.OpenAPIConfig.Info.Title = "spheric-api"
	serverConfig.OpenAPIConfig.Info.Version = "0.1"

	if utilfeature.DefaultFeatureGate.Enabled(features.OpenAPIV3) {
		serverConfig.OpenAPIV3Config = genericapiserver.DefaultOpenAPIV3Config(sphericopenapi.GetOpenAPIDefinitions, openapi.NewDefinitionNamer(api.Scheme))
		serverConfig.OpenAPIV3Config.Info.Title = "spheric-api"
		serverConfig.OpenAPIV3Config.Info.Version = "0.1"
	}

	if err := o.RecommendedOptions.ApplyTo(serverConfig); err != nil {
		return nil, err
	}

	apiResourceConfig := NewResourceConfig()

	config := &apiserver.Config{
		GenericConfig: serverConfig,
		ExtraConfig: apiserver.ExtraConfig{
			APIResourceConfigSource: apiResourceConfig,
			MachinePoolletConfig:    o.MachinePoolletConfig,
		},
	}

	if config.GenericConfig.EgressSelector != nil {
		// Use the config.GenericConfig.EgressSelector lookup to find the dialer to connect to the machinepoollet
		config.ExtraConfig.MachinePoolletConfig.Lookup = config.GenericConfig.EgressSelector.Lookup
	}

	return config, nil
}

func (o *SphericAPIServerOptions) Run(ctx context.Context) error {
	config, err := o.Config()
	if err != nil {
		return err
	}

	server, err := config.Complete().New()
	if err != nil {
		return err
	}

	server.GenericAPIServer.AddPostStartHookOrDie("start-spheric-api-server-informers", func(context genericapiserver.PostStartHookContext) error {
		config.GenericConfig.SharedInformerFactory.Start(context.StopCh)
		o.SharedInformerFactory.Start(context.StopCh)
		return nil
	})

	return server.GenericAPIServer.PrepareRun().Run(ctx.Done())
}
