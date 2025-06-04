// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package clientcmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"k8s.io/client-go/util/homedir"
	"spheric.cloud/spheric/irictl/decoder"
)

const (
	RecommendedConfigPathFlag   = "config"
	RecommendedConfigPathEnvVar = "IRICTL_INSTANCE_CONFIG"
	RecommendedHomeDir          = ".irictl"
	RecommendedFileName         = "config"
)

var (
	RecommendedConfigDir               = filepath.Join(homedir.HomeDir(), RecommendedHomeDir)
	RecommendedHomeFile                = filepath.Join(RecommendedConfigDir, RecommendedFileName)
	RecommendedInstanceRuntimeEndpoint = fmt.Sprintf("unix://%s", filepath.Join("/run", "vee", "vee.sock"))
)

type Column struct {
	Name     string `json:"name"`
	Template string `json:"template"`
}

type Config struct {
}

func DefaultConfig() *Config {
	return &Config{}
}

func ReadConfig(data []byte) (*Config, error) {
	cfg := &Config{}
	if err := decoder.Decode(data, cfg); err != nil {
		return nil, fmt.Errorf("error decoding config: %w", err)
	}
	return cfg, nil
}

func ReadConfigFile(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}
	return ReadConfig(data)
}

func GetConfig(filename string) (*Config, error) {
	if filename != "" {
		return ReadConfigFile(filename)
	}

	if configPath := os.Getenv(RecommendedConfigPathEnvVar); configPath != "" {
		cfg, err := ReadConfigFile(configPath)
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			return nil, err
		}
		if err == nil {
			return cfg, nil
		}
	}

	cfg, err := ReadConfigFile(RecommendedHomeFile)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}
	if err == nil {
		return cfg, nil
	}

	return DefaultConfig(), nil
}
