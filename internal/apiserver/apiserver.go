// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package apiserver

import (
	"fmt"

	"k8s.io/apiserver/pkg/registry/generic"
	genericapiserver "k8s.io/apiserver/pkg/server"
	serverstorage "k8s.io/apiserver/pkg/server/storage"
	ctrl "sigs.k8s.io/controller-runtime"
	corerest "spheric.cloud/spheric/internal/registry/core/rest"
	"spheric.cloud/spheric/internal/spherelet/client"
)

const (
	SphericComponentName = "spheric"
)

var (
	logf = ctrl.Log.WithName("apiserver")
)

// ExtraConfig holds custom apiserver config
type ExtraConfig struct {
	APIResourceConfigSource serverstorage.APIResourceConfigSource
	SphereletConfig         client.SphereletClientConfig
}

// Config defines the config for the apiserver
type Config struct {
	GenericConfig *genericapiserver.RecommendedConfig
	ExtraConfig   ExtraConfig
}

// SphericAPIServer contains state for a Kubernetes cluster master/api server.
type SphericAPIServer struct {
	GenericAPIServer *genericapiserver.GenericAPIServer
}

type completedConfig struct {
	GenericConfig genericapiserver.CompletedConfig
	ExtraConfig   *ExtraConfig
}

// CompletedConfig embeds a private pointer that cannot be instantiated outside of this package.
type CompletedConfig struct {
	*completedConfig
}

// Complete fills in any fields not set that are required to have valid data. It's mutating the receiver.
func (cfg *Config) Complete() CompletedConfig {
	c := completedConfig{
		cfg.GenericConfig.Complete(),
		&cfg.ExtraConfig,
	}

	return CompletedConfig{&c}
}

type RESTStorageProvider interface {
	GroupName() string
	NewRESTStorage(apiResourceConfigSource serverstorage.APIResourceConfigSource, restOptionsGetter generic.RESTOptionsGetter) (genericapiserver.APIGroupInfo, bool, error)
}

// New returns a new instance of SphericAPIServer from the given config.
func (c completedConfig) New() (*SphericAPIServer, error) {
	genericServer, err := c.GenericConfig.New("sample-apiserver", genericapiserver.NewEmptyDelegate())
	if err != nil {
		return nil, err
	}

	s := &SphericAPIServer{
		GenericAPIServer: genericServer,
	}

	apiResourceConfigSource := c.ExtraConfig.APIResourceConfigSource
	restStorageProviders := []RESTStorageProvider{
		corerest.StorageProvider{
			SphereletClientConfig: c.ExtraConfig.SphereletConfig,
		},
	}

	var apiGroupsInfos []*genericapiserver.APIGroupInfo
	for _, restStorageProvider := range restStorageProviders {
		groupName := restStorageProvider.GroupName()
		logf := logf.WithValues("GroupName", groupName)

		if !apiResourceConfigSource.AnyResourceForGroupEnabled(groupName) {
			logf.V(1).Info("Skipping disabled api group")
			continue
		}

		apiGroupInfo, enabled, err := restStorageProvider.NewRESTStorage(apiResourceConfigSource, c.GenericConfig.RESTOptionsGetter)
		if err != nil {
			return nil, fmt.Errorf("error initializing api group %s: %w", groupName, err)
		}
		if !enabled {
			logf.Info("API Group is not enabled, skipping")
			continue
		}

		if postHookProvider, ok := restStorageProvider.(genericapiserver.PostStartHookProvider); ok {
			name, hook, err := postHookProvider.PostStartHook()
			if err != nil {
				return nil, fmt.Errorf("error building post start hook: %w", err)
			}

			if err := s.GenericAPIServer.AddPostStartHook(name, hook); err != nil {
				return nil, err
			}
		}

		apiGroupsInfos = append(apiGroupsInfos, &apiGroupInfo)
	}

	if err := s.GenericAPIServer.InstallAPIGroups(apiGroupsInfos...); err != nil {
		return nil, err
	}

	return s, nil
}
