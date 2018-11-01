// Copyright © 2017 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package agent

import (
	"expvar"
	"net"
	"net/http"
	"os"
	"os/signal"

	"github.com/circonus-labs/circonus-logwatch/internal/config"
	"github.com/circonus-labs/circonus-logwatch/internal/configs"
	"github.com/circonus-labs/circonus-logwatch/internal/metrics/circonus"
	"github.com/circonus-labs/circonus-logwatch/internal/metrics/logonly"
	"github.com/circonus-labs/circonus-logwatch/internal/metrics/statsd"
	"github.com/circonus-labs/circonus-logwatch/internal/release"
	"github.com/circonus-labs/circonus-logwatch/internal/watcher"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	tomb "gopkg.in/tomb.v2"
)

// Agent holds the main circonus-logwatch process
type Agent struct {
	watchers   []*watcher.Watcher
	signalCh   chan os.Signal
	t          tomb.Tomb
	destClient metrics.Destination
	svrHTTP    *http.Server
}

func init() {
	http.Handle("/stats", expvar.Handler())
}

// New returns a new agent instance
func New() (*Agent, error) {
	a := Agent{
		signalCh: make(chan os.Signal),
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
		w, err := watcher.New(a.destClient, cfg)
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
	a.t.Go(a.handleSignals)

	for _, w := range a.watchers {
		a.t.Go(w.Start)
	}

	a.t.Go(a.serveMetrics)

	log.Debug().
		Int("pid", os.Getpid()).
		Str("name", release.NAME).
		Str("ver", release.VERSION).Msg("Starting wait")

	return a.t.Wait()
}

// Stop cleans up and shuts down the Agent
func (a *Agent) Stop() {
	a.stopSignalHandler()

	for _, w := range a.watchers {
		w.Stop()
	}

	a.svrHTTP.Close()

	log.Debug().
		Int("pid", os.Getpid()).
		Str("name", release.NAME).
		Str("ver", release.VERSION).Msg("Stopped")

	if a.t.Alive() {
		a.t.Kill(nil)
	}
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
