// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"os/signal"

	cloudhypervisor "spheric.cloud/spheric/cloud-hypervisor"
)

func run() error {
	ctx, cancel := signal.NotifyContext(context.Background())
	defer cancel()

	slog.Info("Starting API")
	c, err := cloudhypervisor.Default.Start(ctx, "/tmp/cloud-hypervisor.sock")
	if err != nil {
		return err
	}

	slog.Info("Ping...")
	res, err := c.PingVMM(ctx)
	if err != nil {
		return err
	}

	slog.Info("Pong!")
	_ = json.NewEncoder(os.Stdout).Encode(res)
	return nil
}

func main() {
	if err := run(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
