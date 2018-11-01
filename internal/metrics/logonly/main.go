// Copyright © 2017 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package logonly

import (
	"sync"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// LogOnly defines logging metrics only destination
type LogOnly struct {
	logger zerolog.Logger
}

var (
	client *LogOnly
	once   sync.Once
)

// New creates a new log only destination
func New() (*LogOnly, error) {

	once.Do(func() {
		client = &LogOnly{
			logger: log.With().Str("pkg", "dest-log").Logger(),
		}
	})

	return client, nil
}

// Start is a NOP for log only destination
func (c *LogOnly) Start() error {
	// NOP
	return nil
}

// Stop is a NOP for log only destination
func (c *LogOnly) Stop() error {
	// NOP
	return nil
}

// IncrementCounter increments a counter - type 'c'
func (c *LogOnly) IncrementCounter(metric string) error { // counter (monotonically increasing value)
	c.logger.Info().Str("name", metric).Interface("value", 1).Msg("metric")
	return nil
}

// IncrementCounterByValue sends value to add to counter - type 'c'
func (c *LogOnly) IncrementCounterByValue(metric string, value uint64) error { // counter (monotonically increasing value)
	c.logger.Info().Str("name", metric).Interface("value", value).Msg("metric")
	return nil
}

// SetGaugeValue sets a gauge metric to the specified value - type 'g'
func (c *LogOnly) SetGaugeValue(metric string, value interface{}) error { // gauge (ints or floats)
	c.logger.Info().Str("name", metric).Interface("value", value).Msg("metric")
	return nil
}

// SetHistogramValue sets a histogram metric to the specified value - type 'h'
func (c *LogOnly) SetHistogramValue(metric string, value float64) error { // histogram
	c.logger.Info().Str("name", metric).Interface("value", value).Msg("metric")
	return nil
}

// SetTimingValue sets a timing metric to the specified value - type 'ms'
func (c *LogOnly) SetTimingValue(metric string, value float64) error { // histogram
	c.logger.Info().Str("name", metric).Interface("value", value).Msg("metric")
	return nil
}

// AddSetValue adds (or increments the counter) for the specified unique value - type 's'
func (c *LogOnly) AddSetValue(metric, value string) error { // set metric (ala statsd, counts unique values)
	c.logger.Info().Str("name", metric).Interface("value", value).Msg("metric")
	return nil
}

// SetTextValue sets a text metric to the specified value - type 't'
func (c *LogOnly) SetTextValue(metric, value string) error { // text metric
	c.logger.Info().Str("name", metric).Interface("value", value).Msg("metric")
	return nil
}
