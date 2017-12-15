// Copyright © 2017 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package watcher

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	stdlog "log"
	"strconv"
	"strings"
	"time"

	"github.com/circonus-labs/circonus-logwatch/internal/config"
	"github.com/circonus-labs/circonus-logwatch/internal/configs"
	"github.com/circonus-labs/circonus-logwatch/internal/metrics"
	"github.com/hpcloud/tail"
	"github.com/maier/go-appstats"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

// New creates a new watcher instance
func New(metricDest metrics.Destination, logConfig *configs.Config) (*Watcher, error) {
	if metricDest == nil {
		return nil, errors.New("invalid metric destination (nil)")
	}
	if logConfig == nil {
		return nil, errors.New("invalid log config (nil)")
	}

	w := Watcher{
		logger:           log.With().Str("pkg", "watcher").Str("log_id", logConfig.ID).Logger(),
		cfg:              logConfig,
		dest:             metricDest,
		metricLines:      make(chan metricLine, metricLineQueueSize),
		metrics:          make(chan metric, metricQueueSize),
		trace:            viper.GetBool(config.KeyDebugMetric),
		statMatchedLines: logConfig.ID + "_lines_matched",
		statTotalLines:   logConfig.ID + "_lines_total",
	}

	appstats.NewInt(w.statMatchedLines)
	appstats.NewInt(w.statTotalLines)

	return &w, nil
}

// Start the watcher
func (w *Watcher) Start() error {
	w.t.Go(w.save)
	w.t.Go(w.parse)
	w.t.Go(w.process)
	return w.t.Wait()
}

// Stop the watcher
func (w *Watcher) Stop() error {
	w.logger.Info().Msg("stopping")
	if w.t.Alive() {
		w.t.Kill(nil)
	}

	return nil
}

// process opens log and checks log lines for matches
func (w *Watcher) process() error {
	cfg := tail.Config{
		Follow:    true,
		ReOpen:    true,
		MustExist: false,
		Location: &tail.SeekInfo{
			Offset: 0,
			Whence: io.SeekEnd,
		},
		Logger: stdlog.New(ioutil.Discard, "", 0),
	}

	if viper.GetBool(config.KeyDebugTail) {
		cfg.Logger = stdlog.New(w.logger.With().Str("pkg", "tail").Logger(), "", 0)
	}

START_TAIL:
	w.logger.Debug().Msg("starting tail")
	tailer, err := tail.TailFile(w.cfg.LogFile, cfg)
	if err != nil {
		w.logger.Error().Err(err).Msg("starting tailer")
		w.t.Kill(err)
		return err
	}

	w.logger.Debug().Msg("tail started, waiting for lines")
	for {
		select {
		case <-w.t.Dying():
			w.logger.Debug().Msg("w.t dying, stopping tail")
			tailer.Cleanup()
			return nil
		case <-tailer.Dying():
			w.logger.Debug().Err(tailer.Err()).Msg("tailer dying, restarting tailer")
			// there is a not well handled scenario in tail where the inotify watcher
			// is closed while the log reopener is waiting for log file creation events
			tailer.Cleanup()
			goto START_TAIL
		case line := <-tailer.Lines:
			if line == nil {
				_, err := tailer.Tell()
				if err != nil {
					w.logger.Error().Err(err).Msg("nil line w/error")
					if !strings.Contains(err.Error(), "file already closed") {
						w.logger.Debug().Msg("!file already closed error, stopping tail")
						tailer.Cleanup()
						w.t.Kill(err)
						return err
					}
				}
				w.logger.Warn().Msg("nil line, ignoring")
				continue
			}
			appstats.IncrementInt(w.statTotalLines)
			if line.Err != nil {
				w.logger.Error().
					Err(line.Err).
					Str("log_line", line.Text).
					Msg("tail line error -- ignoring line")
				continue
			}
			for id, def := range w.cfg.Metrics {
				if w.trace {
					w.logger.Log().
						Int("metric_id", id).
						Str("metric_match", def.Matcher.String()).
						Str("log_line", line.Text).
						Msg("checking rule")
				}
				matches := def.Matcher.FindAllStringSubmatch(line.Text, -1)
				if matches != nil {
					ml := metricLine{
						line:     line.Text,
						metricID: id,
					}
					if len(def.MatchParts) > 0 {
						m := map[string]string{}
						for i, val := range matches[0] {
							if def.MatchParts[i] != "" {
								m[def.MatchParts[i]] = val
							}
						}
						ml.matches = &m
					}
					w.metricLines <- ml
					// NOTE: do not 'break' on match, a single
					//       line may generate multiple metrics.
				}
			}
		}
	}

	return nil
}

// parse log line to extract metric
func (w *Watcher) parse() error {
	for {
		select {
		case <-w.t.Dying():
			return nil
		case l := <-w.metricLines:
			appstats.IncrementInt(w.statMatchedLines)
			if w.trace {
				w.logger.Log().
					Int("metric_id", l.metricID).
					Str("line", l.line).
					Interface("matches", l.matches).
					Msg("matched, parsing metric line")
			}

			r := w.cfg.Metrics[l.metricID]
			m := metric{
				Name: fmt.Sprintf("%s`%s", w.cfg.ID, r.Name),
				Type: r.Type,
			}

			if m.Type == "c" {
				m.Value = "1" // default to simple incrment by 1
			}

			if l.matches != nil {
				if r.ValueKey != "" {
					v, ok := (*l.matches)[r.ValueKey]
					if !ok {
						w.logger.Warn().
							Str("value_key", r.ValueKey).
							Str("line", l.line).
							Interface("matches", *l.matches).
							Msg("'Value' key defined but not found in matches")
						continue
					}
					m.Value = v
				}
				if r.Namer != nil {
					var b bytes.Buffer
					r.Namer.Execute(&b, *l.matches)
					m.Name = fmt.Sprintf("%s`%s", w.cfg.ID, b.String())
				}
			}

			w.metrics <- m
		}
	}
	return nil
}

// save metrics to configured destination
func (w *Watcher) save() error {
	for {
		select {
		case <-w.t.Dying():
			return nil
		case m := <-w.metrics:
			w.logger.Info().
				Str("metric", fmt.Sprintf("%#v", m)).
				Msg("sending")

			switch m.Type {
			case "c":
				v, err := strconv.ParseUint(m.Value, 10, 64)
				if err != nil {
					w.logger.Warn().Err(err).Msg(m.Name)
				} else {
					w.dest.IncrementCounterByValue(m.Name, v)
				}
			case "g":
				w.dest.SetGaugeValue(m.Name, m.Value)
			case "h":
				v, err := strconv.ParseFloat(m.Value, 64)
				if err != nil {
					w.logger.Warn().Err(err).Msg(m.Name)
				} else {
					w.dest.SetHistogramValue(m.Name, v)
				}
			case "ms":
				// parse as float
				v, errFloat := strconv.ParseFloat(m.Value, 64)
				if errFloat == nil {
					w.dest.SetTimingValue(m.Name, v)
					continue
				}
				// try parsing as a duration (e.g. 60ms, 1m, 3s)
				dur, errDuration := time.ParseDuration(m.Value)
				if errDuration != nil {
					w.logger.Warn().Err(errFloat).Err(errDuration).Str("metric", m.Name).Msg("failed to parse timing as float or duration")
					continue
				}
				w.dest.SetTimingValue(m.Name, float64(dur/time.Millisecond))
			case "s":
				w.dest.AddSetValue(m.Name, m.Value)
			case "t":
				w.dest.SetTextValue(m.Name, m.Value)
			default:
				w.logger.Info().
					Str("type", m.Type).
					Str("name", m.Name).
					Interface("val", m.Value).
					Msg("metric, unknown type")
			}
		}
	}
	return nil
}
