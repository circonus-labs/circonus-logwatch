// Copyright © 2017 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package agent

import (
	"context"
	"errors"
	"expvar"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/circonus-labs/circonus-logwatch/internal/config"
	"github.com/circonus-labs/circonus-logwatch/internal/configs"
	"github.com/circonus-labs/circonus-logwatch/internal/metrics"
	"github.com/circonus-labs/circonus-logwatch/internal/metrics/circonus"
	"github.com/circonus-labs/circonus-logwatch/internal/metrics/logonly"
	"github.com/circonus-labs/circonus-logwatch/internal/metrics/statsd"
	"github.com/circonus-labs/circonus-logwatch/internal/release"
	"github.com/circonus-labs/circonus-logwatch/internal/watcher"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"
)

// Agent holds the main circonus-logwatch process.
type Agent struct {
	groupCtx    context.Context
	destClient  metrics.Destination
	group       *errgroup.Group
	groupCancel context.CancelFunc
	signalCh    chan os.Signal
	svrHTTP     *http.Server
	watchers    []*watcher.Watcher
}

func init() {
	http.Handle("/stats", expvar.Handler())
}

// New returns a new agent instance.
func New() (*Agent, error) {
	ctx, cancel := context.WithCancel(context.Background())
	g, gctx := errgroup.WithContext(ctx)

	a := Agent{
		group:       g,
		groupCtx:    gctx,
		groupCancel: cancel,
		signalCh:    make(chan os.Signal, 10),
	}

	//
	// validate the configuration
	//
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("config validate: %w", err)
	}

	dest := viper.GetString(config.KeyDestType)
	switch dest {
	case "agent":
		fallthrough
	case "check":
		d, err := circonus.New()
		if err != nil {
			return nil, err
		}
		a.destClient = d

	case "statsd":
		d, err := statsd.New()
		if err != nil {
			return nil, err
		}
		a.destClient = d

	case "log":
		d, err := logonly.New()
		if err != nil {
			return nil, err
		}
		a.destClient = d

	default:
		return nil, fmt.Errorf("unknown metric destination (%s)", dest)
	}

	cfgs, err := configs.Load()
	if err != nil {
		return nil, err
	}
	if len(cfgs) == 0 {
		return nil, err
	}

	a.watchers = make([]*watcher.Watcher, len(cfgs))
	for idx, cfg := range cfgs {
		w, err := watcher.New(a.groupCtx, a.destClient, cfg)
		if err != nil {
			log.Error().Err(err).Str("id", cfg.ID).Msg("adding watcher, log will NOT be processed")
		}
		a.watchers[idx] = w
	}

	a.svrHTTP = &http.Server{Addr: net.JoinHostPort("localhost", viper.GetString(config.KeyAppStatPort))}
	a.svrHTTP.SetKeepAlivesEnabled(false)

	a.setupSignalHandler()

	return &a, nil
}

// Start the agent.
func (a *Agent) Start() error {
	a.group.Go(a.handleSignals)
	for _, w := range a.watchers {
		a.group.Go(w.Start)
	}
	a.group.Go(a.serveMetrics)

	log.Debug().
		Int("pid", os.Getpid()).
		Str("name", release.NAME).
		Str("ver", release.VERSION).Msg("Started")

	return a.group.Wait()
}

// Stop cleans up and shuts down the Agent.
func (a *Agent) Stop() {
	log.Debug().
		Int("pid", os.Getpid()).
		Str("name", release.NAME).
		Str("ver", release.VERSION).Msg("Stopping")

	a.stopSignalHandler()
	a.groupCancel()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err := a.svrHTTP.Shutdown(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("closing HTTP server")
	}
}

func (a *Agent) serveMetrics() error {
	log.Debug().Str("url", "http://"+a.svrHTTP.Addr+"/stats").Msg("app stats listener")
	if err := a.svrHTTP.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("http server: %w", err)
		}
	}
	return nil
}

// stopSignalHandler disables the signal handler.
func (a *Agent) stopSignalHandler() {
	signal.Stop(a.signalCh)
	signal.Reset() // so a second ctrl-c will force a kill
}
