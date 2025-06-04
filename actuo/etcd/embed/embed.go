// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0

package embed

import (
	"cmp"
	"context"
	"fmt"
	"net"
	"net/url"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/go-logr/logr"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/server/v3/embed"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
)

type logrCore struct {
	logger logr.Logger
	level  zapcore.Level
}

func newLogrCore(logger logr.Logger, level zapcore.Level) *logrCore {
	return &logrCore{
		logger: logger,
		level:  level,
	}
}

func (c *logrCore) Enabled(level zapcore.Level) bool {
	return level >= c.level
}

func (c *logrCore) With(fields []zapcore.Field) zapcore.Core {
	// Convert zap fields to key-value pairs for logr
	var kvs []interface{}
	for _, f := range fields {
		kvs = append(kvs, f.Key, f.Interface)
	}
	return &logrCore{
		logger: c.logger.WithValues(kvs...),
		level:  c.level,
	}
}

func (c *logrCore) Check(entry zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Enabled(entry.Level) {
		return ce.AddCore(entry, c)
	}
	return ce
}

func (c *logrCore) Sync() error {
	return nil
}

func (c *logrCore) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	// Convert zap fields to key-value pairs for logr
	var kvs []interface{}
	for _, f := range fields {
		kvs = append(kvs, f.Key, f.Interface)
	}

	switch entry.Level {
	case zapcore.DebugLevel:
		c.logger.V(1).Info(entry.Message, kvs...)
	case zapcore.InfoLevel:
		c.logger.Info(entry.Message, kvs...)
	case zapcore.WarnLevel:
		c.logger.Info("WARN: "+entry.Message, kvs...)
	case zapcore.ErrorLevel:
		c.logger.Error(nil, entry.Message, kvs...)
	default:
		c.logger.Info(entry.Message, kvs...)
	}
	return nil
}

func newEtcdConfig(dir string, log logr.Logger) (*embed.Config, error) {
	cfg := embed.NewConfig()
	cfg.UnsafeNoFsync = true

	ports, err := getAvailablePorts(2)
	if err != nil {
		return nil, err
	}
	clientURL := url.URL{Scheme: "http", Host: net.JoinHostPort("localhost", strconv.Itoa(ports[0]))}
	peerURL := url.URL{Scheme: "http", Host: net.JoinHostPort("localhost", strconv.Itoa(ports[1]))}

	cfg.ListenPeerUrls = []url.URL{peerURL}
	cfg.AdvertisePeerUrls = []url.URL{peerURL}
	cfg.ListenClientUrls = []url.URL{clientURL}
	cfg.AdvertiseClientUrls = []url.URL{clientURL}
	cfg.InitialCluster = cfg.InitialClusterFromName(cfg.Name)

	if log.GetSink() == nil {
		log = logr.Discard()
	}

	zapLog := zap.New(newLogrCore(log, zapcore.InfoLevel))
	cfg.ZapLoggerBuilder = embed.NewZapLoggerBuilder(zapLog)

	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, err
	}

	cfg.Dir = dir
	return cfg, nil
}

func clientFor(cfg *embed.Config, e *embed.Etcd) (*clientv3.Client, error) {
	tlsConfig, err := cfg.ClientTLSInfo.ClientConfig()
	if err != nil {
		return nil, err
	}

	return clientv3.New(clientv3.Config{
		TLS:         tlsConfig,
		Endpoints:   e.Server.Cluster().ClientURLs(),
		DialTimeout: 10 * time.Second,
		DialOptions: []grpc.DialOption{},
		//Logger:      zaptest.NewLogger(t, zaptest.Level(zapcore.ErrorLevel)).Named("etcd-client"),
	})
}

type Etcd struct {
	cfg *embed.Config

	startMu sync.Mutex
	started bool
	etcd    *embed.Etcd
}

type Options struct {
	Logger logr.Logger
	Dir    string
}

type WithLogger logr.Logger

func (w WithLogger) ApplyTo(o *Options) {
	o.Logger = logr.Logger(w)
}

type WithDir string

func (w WithDir) ApplyTo(o *Options) {
	o.Dir = string(w)
}

func (o *Options) ApplyOptions(opts []Option) *Options {
	for _, opt := range opts {
		opt.ApplyTo(o)
	}
	return o
}

type Option interface {
	ApplyTo(o *Options)
}

func New(opts ...Option) (*Etcd, error) {
	o := (&Options{}).ApplyOptions(opts)

	dir := cmp.Or(o.Dir, "etcd-data")
	log := o.Logger
	if log.GetSink() == nil {
		log = logr.Discard()
	}

	cfg, err := newEtcdConfig(dir, log)
	if err != nil {
		return nil, err
	}

	return &Etcd{
		cfg: cfg,
	}, nil
}

func (e *Etcd) Start() error {
	e.startMu.Lock()
	defer e.startMu.Unlock()
	if e.started {
		return fmt.Errorf("etcd already started")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	embedEtcd, err := startEtcd(ctx, e.cfg)
	if err != nil {
		return err
	}

	e.etcd = embedEtcd
	e.started = true
	return nil
}

func (e *Etcd) Stop() error {
	e.startMu.Lock()
	defer e.startMu.Unlock()
	if !e.started {
		return fmt.Errorf("etcd not started")
	}

	e.etcd.Close()
	e.started = false
	return nil
}

func (e *Etcd) NewClient() (*clientv3.Client, error) {
	e.startMu.Lock()
	defer e.startMu.Unlock()

	if !e.started {
		return nil, fmt.Errorf("etcd not started")
	}

	return clientFor(e.cfg, e.etcd)
}

func startEtcd(ctx context.Context, cfg *embed.Config) (*embed.Etcd, error) {
	e, err := embed.StartEtcd(cfg)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	select {
	case <-e.Server.ReadyNotify():
	case <-ctx.Done():
		e.Close()
		return nil, ctx.Err()
	}
	return e, nil
}

// getAvailablePort returns a TCP port that is available for binding.
func getAvailablePorts(count int) ([]int, error) {
	ports := []int{}
	for i := 0; i < count; i++ {
		l, err := net.Listen("tcp", ":0")
		if err != nil {
			return nil, fmt.Errorf("could not bind to a port: %v", err)
		}
		// It is possible but unlikely that someone else will bind this port before we get a chance to use it.
		defer l.Close()
		ports = append(ports, l.Addr().(*net.TCPAddr).Port)
	}
	return ports, nil
}
