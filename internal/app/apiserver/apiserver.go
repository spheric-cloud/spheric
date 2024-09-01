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

	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/version"
	utilversion "k8s.io/apiserver/pkg/util/version"
	"k8s.io/component-base/featuregate"
	baseversion "k8s.io/component-base/version"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apiserver/pkg/admission"
	"k8s.io/apiserver/pkg/endpoints/openapi"
	genericapiserver "k8s.io/apiserver/pkg/server"
	genericoptions "k8s.io/apiserver/pkg/server/options"
	serverstorage "k8s.io/apiserver/pkg/server/storage"
	utilfeature "k8s.io/apiserver/pkg/util/feature"
	netutils "k8s.io/utils/net"
	corev1alpha1 "spheric.cloud/spheric/api/core/v1alpha1"
	"spheric.cloud/spheric/client-go/informers"
	sphericopenapi "spheric.cloud/spheric/client-go/openapi"
	clientset "spheric.cloud/spheric/client-go/spheric"
	sphericinitializer "spheric.cloud/spheric/internal/admission/initializer"
	"spheric.cloud/spheric/internal/admission/plugin/instancediskdevices"
	"spheric.cloud/spheric/internal/api"
	"spheric.cloud/spheric/internal/apis/core"
	"spheric.cloud/spheric/internal/apiserver"
	"spheric.cloud/spheric/internal/spherelet/client"
)

const defaultEtcdPathPrefix = "/registry/spheric.cloud"

func NewResourceConfig() *serverstorage.ResourceConfig {
	cfg := serverstorage.NewResourceConfig()
	cfg.EnableVersions(
		corev1alpha1.SchemeGroupVersion,
	)
	return cfg
}

func SphericVersionToKubeVersion(ver *version.Version) *version.Version {
	if ver.Major() != 1 {
		return nil
	}
	kubeVer := utilversion.DefaultKubeEffectiveVersion().BinaryVersion()
	// "1.2" maps to kubeVer
	offset := int(ver.Minor()) - 2
	mappedVer := kubeVer.OffsetMinor(offset)
	if mappedVer.GreaterThan(kubeVer) {
		return kubeVer
	}
	return mappedVer
}

type SphericAPIServerOptions struct {
	RecommendedOptions *genericoptions.RecommendedOptions
	SphereletConfig    client.SphereletClientConfig

	SharedInformerFactory informers.SharedInformerFactory
}

func (o *SphericAPIServerOptions) AddFlags(fs *pflag.FlagSet) {
	o.RecommendedOptions.AddFlags(fs)

	// spherelet related flags:
	fs.StringSliceVar(&o.SphereletConfig.PreferredAddressTypes, "spherelet-preferred-address-types", o.SphereletConfig.PreferredAddressTypes,
		"List of the preferred FleetAddressTypes to use for spherelet connections.")

	fs.DurationVar(&o.SphereletConfig.HTTPTimeout, "spherelet-timeout", o.SphereletConfig.HTTPTimeout,
		"Timeout for spherelet operations.")

	fs.StringVar(&o.SphereletConfig.CertFile, "spherelet-client-certificate", o.SphereletConfig.CertFile,
		"Path to a client cert file for TLS.")

	fs.StringVar(&o.SphereletConfig.KeyFile, "spherelet-client-key", o.SphereletConfig.KeyFile,
		"Path to a client key file for TLS.")

	fs.StringVar(&o.SphereletConfig.CAFile, "spherelet-certificate-authority", o.SphereletConfig.CAFile,
		"Path to a cert file for the certificate authority.")
}

func NewSphericAPIServerOptions() *SphericAPIServerOptions {
	o := &SphericAPIServerOptions{
		RecommendedOptions: genericoptions.NewRecommendedOptions(
			defaultEtcdPathPrefix,
			api.Codecs.LegacyCodec(
				corev1alpha1.SchemeGroupVersion,
			),
		),
		SphereletConfig: client.SphereletClientConfig{
			Port:         12319,
			ReadOnlyPort: 12320,
			PreferredAddressTypes: []string{
				string(core.FleetHostName),

				// internal, preferring DNS if reported
				string(core.FleetInternalDNS),
				string(core.FleetInternalIP),

				// external, preferring DNS if reported
				string(core.FleetExternalDNS),
				string(core.FleetExternalIP),
			},
			HTTPTimeout: time.Duration(5) * time.Second,
		},
	}
	o.RecommendedOptions.Etcd.StorageConfig.EncodeVersioner = runtime.NewMultiGroupVersioner(
		corev1alpha1.SchemeGroupVersion,
		schema.GroupKind{Group: corev1alpha1.GroupName},
	)
	return o
}

func NewCommandStartSphericAPIServer(ctx context.Context, defaults *SphericAPIServerOptions) *cobra.Command {
	o := *defaults
	cmd := &cobra.Command{
		Short: "Launch a spheric API server",
		Long:  "Launch a spheric API server",
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
	cmd.SetContext(ctx)

	flags := cmd.Flags()
	o.AddFlags(flags)

	utilversion.DefaultComponentGlobalsRegistry.AddFlags(flags)

	defaultSphericVersion := "1.2"

	_, sphericFeatureGate := utilversion.DefaultComponentGlobalsRegistry.ComponentGlobalsOrRegister(
		apiserver.SphericComponentName, utilversion.NewEffectiveVersion(defaultSphericVersion),
		featuregate.NewVersionedFeatureGate(version.MustParse(defaultSphericVersion)))

	utilruntime.Must(sphericFeatureGate.AddVersioned(map[featuregate.Feature]featuregate.VersionedSpecs{}))

	_, _ = utilversion.DefaultComponentGlobalsRegistry.ComponentGlobalsOrRegister(utilversion.DefaultKubeComponent,
		utilversion.NewEffectiveVersion(baseversion.DefaultKubeBinaryVersion), utilfeature.DefaultMutableFeatureGate)

	utilruntime.Must(utilversion.DefaultComponentGlobalsRegistry.SetEmulationVersionMapping(apiserver.SphericComponentName, utilversion.DefaultKubeComponent, SphericVersionToKubeVersion))

	return cmd
}

func (o *SphericAPIServerOptions) Validate(args []string) error {
	var errors []error
	errors = append(errors, o.RecommendedOptions.Validate()...)
	errors = append(errors, utilfeature.DefaultMutableFeatureGate.Validate()...)
	errors = append(errors, utilversion.DefaultComponentGlobalsRegistry.Validate()...)
	return utilerrors.NewAggregate(errors)
}

func (o *SphericAPIServerOptions) Complete() error {
	instancediskdevices.Register(o.RecommendedOptions.Admission.Plugins)

	o.RecommendedOptions.Admission.RecommendedPluginOrder = append(
		o.RecommendedOptions.Admission.RecommendedPluginOrder,
		instancediskdevices.PluginName,
	)

	return nil
}

func (o *SphericAPIServerOptions) Config() (*apiserver.Config, error) {
	if err := o.RecommendedOptions.SecureServing.MaybeDefaultWithSelfSignedCerts("localhost", nil, []net.IP{netutils.ParseIPSloppy("127.0.0.1")}); err != nil {
		return nil, fmt.Errorf("error creating self-signed certificates: %w", err)
	}

	o.RecommendedOptions.ExtraAdmissionInitializers = func(c *genericapiserver.RecommendedConfig) ([]admission.PluginInitializer, error) {
		sphericClient, err := clientset.NewForConfig(c.LoopbackClientConfig)
		if err != nil {
			return nil, err
		}

		informerFactory := informers.NewSharedInformerFactory(sphericClient, c.LoopbackClientConfig.Timeout)
		o.SharedInformerFactory = informerFactory

		genericInitializer := sphericinitializer.New(sphericClient, informerFactory)

		return []admission.PluginInitializer{
			genericInitializer,
		}, nil
	}

	serverConfig := genericapiserver.NewRecommendedConfig(api.Codecs)

	serverConfig.OpenAPIConfig = genericapiserver.DefaultOpenAPIConfig(sphericopenapi.GetOpenAPIDefinitions, openapi.NewDefinitionNamer(api.Scheme))
	serverConfig.OpenAPIConfig.Info.Title = "spheric-api"
	serverConfig.OpenAPIConfig.Info.Version = "0.1"

	serverConfig.OpenAPIV3Config = genericapiserver.DefaultOpenAPIV3Config(sphericopenapi.GetOpenAPIDefinitions, openapi.NewDefinitionNamer(api.Scheme))
	serverConfig.OpenAPIV3Config.Info.Title = "spheric-api"
	serverConfig.OpenAPIV3Config.Info.Version = "0.1"

	serverConfig.FeatureGate = utilversion.DefaultComponentGlobalsRegistry.FeatureGateFor(utilversion.DefaultKubeComponent)
	serverConfig.EffectiveVersion = utilversion.DefaultComponentGlobalsRegistry.EffectiveVersionFor(apiserver.SphericComponentName)

	if err := o.RecommendedOptions.ApplyTo(serverConfig); err != nil {
		return nil, err
	}

	apiResourceConfig := NewResourceConfig()

	config := &apiserver.Config{
		GenericConfig: serverConfig,
		ExtraConfig: apiserver.ExtraConfig{
			APIResourceConfigSource: apiResourceConfig,
			SphereletConfig:         o.SphereletConfig,
		},
	}

	if config.GenericConfig.EgressSelector != nil {
		// Use the config.GenericConfig.EgressSelector lookup to find the dialer to connect to the spherelet
		config.ExtraConfig.SphereletConfig.Lookup = config.GenericConfig.EgressSelector.Lookup
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
		config.GenericConfig.SharedInformerFactory.Start(context.Done())
		o.SharedInformerFactory.Start(context.Done())
		return nil
	})

	return server.GenericAPIServer.PrepareRun().RunWithContext(ctx)
}
