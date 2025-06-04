// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0

package server

import (
	"context"
	"net"
	"net/http"
	"time"
)

type Server struct {
	filename            string
	handler             http.Handler
	shutdownGracePeriod time.Duration
}

func NewServer(filename string, handler http.Handler) *Server {
	return &Server{
		filename:            filename,
		handler:             handler,
		shutdownGracePeriod: 3 * time.Second,
	}
}

func (s *Server) ListenAndServe(ctx context.Context) error {
	ln, err := net.Listen("unix", s.filename)
	if err != nil {
		return err
	}
	defer func() { _ = ln.Close() }()

	var (
		httpSrv = &http.Server{
			Handler: s.handler,
		}
	)
	go func() {
		_ = httpSrv.Serve(ln)
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), s.shutdownGracePeriod)
	defer cancel()
	return httpSrv.Shutdown(shutdownCtx)
}
