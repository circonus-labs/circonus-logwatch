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
	"strings"
	"time"

	cgm "github.com/circonus-labs/circonus-gometrics/v3"
	"github.com/circonus-labs/circonus-logwatch/internal/config"
	"github.com/circonus-labs/circonus-logwatch/internal/config/defaults"
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

		cmc := &cgm.Config{}
		if viper.GetBool(config.KeyDebugCGM) {
			cmc.Debug = true
			cmc.Log = stdlog.New(log.With().Str("pkg", "dest-check").Logger(), "", 0)
		}

		cmc.CheckManager.Check.SubmissionURL = sURL

		interval := viper.GetString(config.KeyDestCfgAgentInterval)
		if interval == "" {
			interval = defaults.AgentInterval
		}
		_, err := time.ParseDuration(interval)
		if err != nil {
			return nil, errors.Wrap(err, "parsing destination interval")
		}
		cmc.Interval = interval

		c, err := cgm.New(cmc)
		if err != nil {
			return nil, errors.Wrap(err, "creating client for destination 'agent'")
		}
		client = c

	case "check":
		cmc := &cgm.Config{}
		if viper.GetBool(config.KeyDebugCGM) {
			cmc.Debug = true
			cmc.Log = stdlog.New(log.With().Str("pkg", "dest-check").Logger(), "", 0)
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

// convert []string to cgm.Tags
func (c *Circonus) tagsToCgmTags(tags []string) cgm.Tags {
	var tagList cgm.Tags
	for _, tag := range tags {
		tp := strings.SplitN(tag, ":", 2)
		if len(tp) != 2 {
			c.logger.Warn().Str("tag", tag).Msg("invalid tag")
			continue
		}
		tagList = append(tagList, cgm.Tag{Category: tp[0], Value: tp[1]})
	}
	return tagList
}

// SetGaugeValue sends a gauge metric
func (c *Circonus) SetGaugeValue(metric string, value interface{}) error { // gauge (ints or floats)
	c.client.Gauge(metric, value)
	return nil
}

// SetGaugeValueWithTags sends a gauge metric with tags
func (c *Circonus) SetGaugeValueWithTags(metric string, tags []string, value interface{}) error { // gauge (ints or floats)
	c.client.GaugeWithTags(metric, c.tagsToCgmTags(tags), value)
	return nil
}

// SetTimingValue sends a timing metric
func (c *Circonus) SetTimingValue(metric string, value float64) error { // histogram
	return c.SetHistogramValue(metric, value)
}

// SetTimingValueWithTags sends a timing metric with tags
func (c *Circonus) SetTimingValueWithTags(metric string, tags []string, value float64) error { // histogram
	return c.SetHistogramValueWithTags(metric, tags, value)
}

// SetHistogramValue sends a histogram metric
func (c *Circonus) SetHistogramValue(metric string, value float64) error { // histogram
	c.client.RecordValue(metric, value)
	return nil
}

// SetHistogramValueWithTags sends a histogram metric with tags
func (c *Circonus) SetHistogramValueWithTags(metric string, tags []string, value float64) error { // histogram
	c.client.RecordValueWithTags(metric, c.tagsToCgmTags(tags), value)
	return nil
}

// IncrementCounter sends a counter increment
func (c *Circonus) IncrementCounter(metric string) error { // counter (monotonically increasing value)
	return c.IncrementCounterByValue(metric, 1)
}

// IncrementCounterWithTags sends a counter increment with tags
func (c *Circonus) IncrementCounterWithTags(metric string, tags []string) error { // counter (monotonically increasing value)
	return c.IncrementCounterByValueWithTags(metric, tags, 1)
}

// IncrementCounterByValue sends value to add to counter
func (c *Circonus) IncrementCounterByValue(metric string, value uint64) error { // counter (monotonically increasing value)
	c.client.IncrementByValue(metric, value)
	return nil
}

// IncrementCounterByValueWithTags sends value to add to counter with tags
func (c *Circonus) IncrementCounterByValueWithTags(metric string, tags []string, value uint64) error { // counter (monotonically increasing value)
	c.client.IncrementByValueWithTags(metric, c.tagsToCgmTags(tags), value)
	return nil
}

// AddSetValue sends a unique value to the set metric
func (c *Circonus) AddSetValue(metric string, value string) error { // set metric (ala statsd, counts unique values)
	_ = c.IncrementCounter(fmt.Sprintf("%s`%s", metric, value))
	return nil
}

// AddSetValueWithTags sends a unique value to the set metric with tags
func (c *Circonus) AddSetValueWithTags(metric string, tags []string, value string) error { // set metric (ala statsd, counts unique values)
	_ = c.IncrementCounterWithTags(fmt.Sprintf("%s`%s", metric, value), tags)
	return nil
}

// SetTextValue sends a text metric
func (c *Circonus) SetTextValue(metric string, value string) error { // text metric
	c.client.SetTextValue(metric, value)
	return nil
}

// SetTextValueWithTags sends a text metric with tags
func (c *Circonus) SetTextValueWithTags(metric string, tags []string, value string) error { // text metric
	c.client.SetTextValueWithTags(metric, c.tagsToCgmTags(tags), value)
	return nil
}
