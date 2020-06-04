// Copyright © 2017 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package circonus

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	stdlog "log"
	"time"

	cgm "github.com/circonus-labs/circonus-gometrics/v3"
	"github.com/circonus-labs/circonus-logwatch/internal/config"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

// Circonus defines an instance of the circonus metrics destination
type Circonus struct {
	logger zerolog.Logger
	client *cgm.CirconusMetrics
}

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

		interval := viper.GetString(config.KeyDestInterval)
		if interval == "" {
			interval = "60s"
		}
		_, err := time.ParseDuration(interval)
		if err != nil {
			return nil, errors.Wrap(err, "parsing destination interval")
		}
		cmc.Interval = interval

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
		cmc.CheckManager.API.TokenKey = viper.GetString(config.KeyAPITokenKey)
		if viper.GetString(config.KeyAPITokenApp) != "" {
			cmc.CheckManager.API.TokenApp = viper.GetString(config.KeyAPITokenApp)
		}
		if viper.GetString(config.KeyAPIURL) != "" {
			cmc.CheckManager.API.URL = viper.GetString(config.KeyAPIURL)
		}
		if file := viper.GetString(config.KeyAPICAFile); file != "" {
			cert, err := ioutil.ReadFile(file)
			if err != nil {
				return nil, errors.Wrapf(err, "reading specified API CA file (%s)", file)
			}
			cp := x509.NewCertPool()
			if !cp.AppendCertsFromPEM(cert) {
				return nil, errors.New("unable to add API CA Certificate to x509 cert pool")
			}
			cmc.CheckManager.API.TLSConfig = &tls.Config{RootCAs: cp}
		}
		if viper.GetString(config.KeyDestCfgCID) != "" {
			cmc.CheckManager.Check.ID = viper.GetString(config.KeyDestCfgCID)
		}
		if viper.GetString(config.KeyDestCfgURL) != "" {
			cmc.CheckManager.Check.SubmissionURL = viper.GetString(config.KeyDestCfgURL)
		}
		if viper.GetString(config.KeyDestCfgSearchTag) != "" {
			cmc.CheckManager.Check.SearchTag = viper.GetString(config.KeyDestCfgSearchTag)
		}
		if viper.GetString(config.KeyDestCfgTarget) != "" {
			cmc.CheckManager.Check.TargetHost = viper.GetString(config.KeyDestCfgTarget)
		}
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
