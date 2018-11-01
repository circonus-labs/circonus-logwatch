// Copyright © 2017 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package agent

import (
	"context"
	"expvar"
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
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"
)

// Agent holds the main circonus-logwatch process
type Agent struct {
	group       *errgroup.Group
	groupCtx    context.Context
	groupCancel context.CancelFunc
	watchers    []*watcher.Watcher
	signalCh    chan os.Signal
	destClient  metrics.Destination
	svrHTTP     *http.Server
}

func init() {
	http.Handle("/stats", expvar.Handler())
}

// New returns a new agent instance
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
		return nil, err
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
		return nil, errors.Errorf("unknown metric destination (%s)", dest)
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

// Start the agent
func (a *Agent) Start() error {
	a.group.Go(a.handleSignals)
	for _, w := range a.watchers {
		a.group.Go(w.Start)
	}
	a.group.Go(a.serveMetrics)

	log.Debug().
		Int("pid", os.Getpid()).
		Str("name", release.NAME).
		Str("ver", release.VERSION).Msg("Starting wait")

	return a.group.Wait()
}

// Stop cleans up and shuts down the Agent
func (a *Agent) Stop() {
	a.stopSignalHandler()
	a.groupCancel()

	// for _, w := range a.watchers {
	// 	w.Stop()
	// }

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err := a.svrHTTP.Shutdown(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("closing HTTP server")
	}

	log.Debug().
		Int("pid", os.Getpid()).
		Str("name", release.NAME).
		Str("ver", release.VERSION).Msg("Stopped")
}

func (a *Agent) serveMetrics() error {
	log.Debug().Str("url", "http://"+a.svrHTTP.Addr+"/stats").Msg("app stats listener")
	if err := a.svrHTTP.ListenAndServe(); err != nil {
		if err != http.ErrServerClosed {
			return errors.Wrap(err, "HTTP server")
		}
	}
	return nil
}

// stopSignalHandler disables the signal handler
func (a *Agent) stopSignalHandler() {
	signal.Stop(a.signalCh)
	signal.Reset() // so a second ctrl-c will force a kill
}
