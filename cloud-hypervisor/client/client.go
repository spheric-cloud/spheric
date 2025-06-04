// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"

	oapiclient "spheric.cloud/spheric/cloud-hypervisor/oapi-client"
)

type Error struct {
	StatusCode int
	Body       []byte
}

func (e *Error) Error() string {
	if e.Body != nil {
		return fmt.Sprintf("unexpected status %d, body %s", e.StatusCode, string(e.Body))
	}
	return fmt.Sprintf("unexpected status %d, no body", e.StatusCode)
}

func IsStatusError(err error, status int) bool {
	var statusErr *Error
	return errors.As(err, &statusErr) && statusErr.StatusCode == status
}

func checkStatus(statusCode int, accepts ...func(statusCode int) bool) error {
	for _, accept := range accepts {
		if accept(statusCode) {
			return nil
		}
	}
	return &Error{statusCode, nil}
}

func checkStatusWithBody(statusCode int, body []byte, accepts ...func(statusCode int) bool) error {
	for _, accept := range accepts {
		if accept(statusCode) {
			return nil
		}
	}
	return &Error{statusCode, body}
}

func isOK(statusCode int) bool {
	return statusCode == http.StatusOK
}

func hasStatus(expected int) func(statusCode int) bool {
	return func(statusCode int) bool {
		return statusCode == expected
	}
}

type Client interface {
	PingVMM(ctx context.Context) (*oapiclient.VmmPingResponse, error)
	CreateVM(ctx context.Context, req oapiclient.CreateVMJSONRequestBody) error
	GetVMInfo(ctx context.Context) (*oapiclient.VmInfo, error)
}

type client struct {
	oapiClient oapiclient.ClientWithResponsesInterface
}

func Connect(socket string) (Client, error) {
	oapiClient, err := oapiclient.NewClientWithResponses("http://localhost/api/v1",
		oapiclient.WithHTTPClient(&http.Client{
			Transport: &http.Transport{
				DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
					return (&net.Dialer{}).DialContext(ctx, "unix", socket)
				},
			},
		}),
	)
	if err != nil {
		return nil, err
	}

	return &client{
		oapiClient: oapiClient,
	}, nil
}

func (c *client) PingVMM(ctx context.Context) (*oapiclient.VmmPingResponse, error) {
	res, err := c.oapiClient.GetVmmPingWithResponse(ctx)
	if err != nil {
		return nil, err
	}
	if err := checkStatusWithBody(res.StatusCode(), res.Body, isOK); err != nil {
		return nil, err
	}
	if res.JSON200 == nil {
		return nil, fmt.Errorf("empty response, body: %s", string(res.Body))
	}
	return res.JSON200, nil
}

func (c *client) CreateVM(ctx context.Context, req oapiclient.CreateVMJSONRequestBody) error {
	res, err := c.oapiClient.CreateVMWithResponse(ctx, req)
	if err != nil {
		return err
	}
	if err := checkStatusWithBody(res.StatusCode(), res.Body, hasStatus(http.StatusNoContent)); err != nil {
		return err
	}
	return nil
}

func (c *client) GetVMInfo(ctx context.Context) (*oapiclient.VmInfo, error) {
	res, err := c.oapiClient.GetVmInfoWithResponse(ctx)
	if err != nil {
		return nil, err
	}
	if err := checkStatusWithBody(res.StatusCode(), res.Body, hasStatus(http.StatusOK)); err != nil {
		return nil, err
	}
	if res.JSON200 == nil {
		return nil, fmt.Errorf("empty response, body: %s", string(res.Body))
	}
	return res.JSON200, nil
}
