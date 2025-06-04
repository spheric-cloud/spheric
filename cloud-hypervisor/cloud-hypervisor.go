// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0

package cloudhypervisor

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"

	"spheric.cloud/spheric/cloud-hypervisor/client"
)

type CloudHypervisor struct {
	command string
}

func NewCloudHypervisor(command string) *CloudHypervisor {
	return &CloudHypervisor{
		command: command,
	}
}

var Default = NewCloudHypervisor("cloud-hypervisor")

func (c *CloudHypervisor) Start(ctx context.Context, apiSocket string) (client.Client, error) {
	cmd := exec.Command(c.command, "--api-socket", apiSocket)

	runErr := make(chan error, 1)
	go func() {
		defer close(runErr)
		res, err := cmd.CombinedOutput()
		if err != nil {
			runErr <- fmt.Errorf("error running api command: %w, output: %s", err, string(res))
		} else {
			runErr <- nil
		}
	}()

	interruptAndWait := func() {
		_ = cmd.Process.Signal(os.Interrupt)
		_, _ = cmd.Process.Wait()
	}

	cl, err := client.Connect(apiSocket)
	if err != nil {
		interruptAndWait()
		return nil, fmt.Errorf("error creating client: %w", err)
	}

	var lastErr error
	check := func() bool {
		_, err := cl.PingVMM(ctx)
		lastErr = err
		return err == nil
	}
	if check() {
		return cl, nil
	}

	t := time.NewTicker(1 * time.Second)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			if lastErr != nil {
				return nil, fmt.Errorf("%w, last error: %v", ctx.Err(), lastErr)
			}
			return nil, ctx.Err()
		case <-t.C:
			if !check() {
				continue
			}
			return cl, nil
		case err := <-runErr:
			if err == nil {
				return nil, fmt.Errorf("cloud-hypervisor exited early without error")
			}
			return nil, err
		}
	}
}
