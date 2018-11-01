// Copyright © 2017 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package statsd

import (
	crand "crypto/rand"
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"net"
	"sync"
	"time"

	"github.com/circonus-labs/circonus-logwatch/internal/config"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

// Statsd defines the relevant properties of a StatsD connection.
type Statsd struct {
	id     string
	port   string
	prefix string
	conn   net.Conn
	logger zerolog.Logger
}

var (
	client *Statsd
	once   sync.Once
)

func init() {
	n, err := crand.Int(crand.Reader, big.NewInt(math.MaxInt64))
	if err != nil {
		rand.Seed(time.Now().UTC().UnixNano())
		return
	}
	rand.Seed(n.Int64())
}

// New initializes a udp connection to the localhost
func New() (*Statsd, error) {
	id := viper.GetString(config.KeyDestCfgID)
	if id == "" {
		return nil, errors.Errorf("invalid id, empty")
	}

	port := viper.GetString(config.KeyDestCfgPort)
	if port == "" {
		return nil, errors.Errorf("invalid port, empty")
	}

	once.Do(func() {
		client = &Statsd{
			id:     id,
			port:   port,
			prefix: viper.GetString(config.KeyDestCfgStatsdPrefix) + id + "`",
			logger: log.With().Str("pkg", "dest-statsd").Logger(),
		}
	})

	return client, nil
}

// Start the statsd Statsd
func (c *Statsd) Start() error {
	return c.open()
}

// Stop the statsd Statsd
func (c *Statsd) Stop() error {
	if c.conn == nil {
		return nil
	}
	return c.conn.Close()
}

// SetGaugeValue sends a gauge metric
func (c *Statsd) SetGaugeValue(metric string, value interface{}) error { // gauge (ints or floats)
	v, err := getGaugeValue(value)
	if err != nil {
		return err
	}
	return c.send(fmt.Sprintf("%s:%s|g", metric, v))
}

// SetTimingValue sends a timing metric
func (c *Statsd) SetTimingValue(metric string, value float64) error { // histogram
	return c.SetHistogramValue(metric, value)
}

// SetHistogramValue sends a histogram metric
func (c *Statsd) SetHistogramValue(metric string, value float64) error { // histogram
	return c.send(fmt.Sprintf("%s:%e|ms", metric, value))
}

// IncrementCounter sends a counter increment
func (c *Statsd) IncrementCounter(metric string) error { // counter (monotonically increasing value)
	return c.IncrementCounterByValue(metric, 1)
}

// IncrementCounterByValue sends value to add to counter
func (c *Statsd) IncrementCounterByValue(metric string, value uint64) error { // counter (monotonically increasing value)
	return c.send(fmt.Sprintf("%s:%d|c", metric, value))
}

// AddSetValue sends a unique value to the set metric
func (c *Statsd) AddSetValue(metric string, value string) error { // set metric (ala statsd, counts unique values)
	return c.send(fmt.Sprintf("%s:%s|s", metric, value))
}

// SetTextValue sends a text metric
func (c *Statsd) SetTextValue(metric string, value string) error { // text metric
	return c.send(fmt.Sprintf("%s:%s|t", metric, value))
}

// send stats data to udp statsd daemon
//
// Outgoing metric format:
//
//   name:value|type[@rate]
//
// e.g.
//   foo:1|c
//   foo:1|c@0.5
//   bar:2.5|ms
//   bar:2.5|ms@.25
//   baz:25|g
//   qux:abcd123|s
//   dib:38.282|h
//   dab:yadda yadda yadda|t
func (c *Statsd) send(metric string) error {
	if c.conn == nil {
		if err := c.open(); err != nil {
			return err
		}
	}

	m := c.prefix + metric

	c.logger.Debug().Str("metric", m).Msg("sending")

	_, err := fmt.Fprintf(c.conn, m)
	if err != nil {
		return err
	}

	return nil
}

// open udp connection
func (c *Statsd) open() error {
	if c.conn != nil {
		c.conn.Close()
	}

	conn, err := net.Dial("udp", net.JoinHostPort("", c.port))
	if err != nil {
		return err
	}

	c.conn = conn

	return nil
}

// getGaugeValue as string from interface
func getGaugeValue(value interface{}) (string, error) {
	vs := ""
	switch v := value.(type) {
	case int:
		vs = fmt.Sprintf("%d", v)
	case int8:
		vs = fmt.Sprintf("%d", v)
	case int16:
		vs = fmt.Sprintf("%d", v)
	case int32:
		vs = fmt.Sprintf("%d", v)
	case int64:
		vs = fmt.Sprintf("%d", v)
	case uint:
		vs = fmt.Sprintf("%d", v)
	case uint8:
		vs = fmt.Sprintf("%d", v)
	case uint16:
		vs = fmt.Sprintf("%d", v)
	case uint32:
		vs = fmt.Sprintf("%d", v)
	case uint64:
		vs = fmt.Sprintf("%d", v)
	case float32:
		vs = fmt.Sprintf("%f", v)
	case float64:
		vs = fmt.Sprintf("%f", v)
	default:
		return "", errors.Errorf("unknown type for value %v", v)
	}
	return vs, nil
}
