// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0

package version

import "runtime/debug"

var (
	Commit string

	Version = "unknown"
)

const (
	RuntimeName = "vee"
)

func init() {
	info, ok := debug.ReadBuildInfo()
	if !ok || info == nil {
		return
	}

	if v := info.Main.Version; v != "" {
		Version = v
	}

	for _, setting := range info.Settings {
		if setting.Key == "vcs.revision" {
			Commit = setting.Value
		}
	}
}
