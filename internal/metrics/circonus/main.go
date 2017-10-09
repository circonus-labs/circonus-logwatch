// Copyright © 2017 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package circonus

import (
	"fmt"
	stdlog "log"

	cgm "github.com/circonus-labs/circonus-gometrics"
	"github.com/circonus-labs/circonus-logwatch/internal/config"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

// New returns a new instance of the circonus metrics destination
func New() (*Circonus, error) {
	var client *cgm.CirconusMetrics

	logger := log.With().Str("pkg", "circonus").Logger()
	dest := viper.GetString(config.KeyDestType)

	switch dest {
	case "agent":
		sURL := viper.GetString(config.KeyDestAgentURL)
		if sURL == "" {
			return nil, errors.Errorf("invalid agent url defined (empty)")
		}
		cmc := &cgm.Config{
			Debug: viper.GetBool(config.KeyDebugCGM),
			Log:   stdlog.New(log.With().Str("pkg", "dest-agent").Logger(), "", 0),
		}
		cmc.Interval = "60s"
		cmc.CheckManager.Check.SubmissionURL = sURL
		c, err := cgm.New(cmc)
		if err != nil {
			return nil, errors.Wrap(err, "creating client for destination 'agent'")
		}
		client = c

	case "check":
		cmc := &cgm.Config{
			Debug: viper.GetBool(config.KeyDebugCGM),
			Log:   stdlog.New(log.With().Str("pkg", "dest-check").Logger(), "", 0),
		}
		cmc.CheckManager.Check.ID = viper.GetString(config.KeyDestCfgCID)
		cmc.CheckManager.Check.SubmissionURL = viper.GetString(config.KeyDestCfgURL)
		cmc.CheckManager.Check.SearchTag = viper.GetString(config.KeyDestCfgSearchTag)
		cmc.CheckManager.Check.TargetHost = viper.GetString(config.KeyDestCfgTarget)
		c, err := cgm.New(cmc)
		if err != nil {
			return nil, errors.Wrap(err, "creating client for destination 'check'")
		}
		client = c

	default:
		return nil, errors.Errorf("unknown destination type for circonus client %s", dest)
	}

	return &Circonus{client: client, logger: logger}, nil
}

// Start NOOP cgm starts as it is initialized
func (c *Circonus) Start() error {
	// noop
	return nil
}

// Stop flushes any outstanding metrics
func (c *Circonus) Stop() error {
	c.client.Flush()
	return nil
}

// SetGaugeValue sends a gauge metric
func (c *Circonus) SetGaugeValue(metric string, value interface{}) error { // gauge (ints or floats)
	c.client.Gauge(metric, value)
	return nil
}

// SetTimingValue sends a timing metric
func (c *Circonus) SetTimingValue(metric string, value float64) error { // histogram
	return c.SetHistogramValue(metric, value)
}

// SetHistogramValue sends a histogram metric
func (c *Circonus) SetHistogramValue(metric string, value float64) error { // histogram
	c.client.RecordValue(metric, value)
	return nil
}

// IncrementCounter sends a counter increment
func (c *Circonus) IncrementCounter(metric string) error { // counter (monotonically increasing value)
	return c.IncrementCounterByValue(metric, 1)
}

// IncrementCounterByValue sends value to add to counter
func (c *Circonus) IncrementCounterByValue(metric string, value uint64) error { // counter (monotonically increasing value)
	c.client.IncrementByValue(metric, value)
	return nil
}

// AddSetValue sends a unique value to the set metric
func (c *Circonus) AddSetValue(metric string, value string) error { // set metric (ala statsd, counts unique values)
	c.IncrementCounter(fmt.Sprintf("%s`%s", metric, value))
	return nil
}

// SetTextValue sends a text metric
func (c *Circonus) SetTextValue(metric string, value string) error { // text metric
	c.client.SetTextValue(metric, value)
	return nil
}
