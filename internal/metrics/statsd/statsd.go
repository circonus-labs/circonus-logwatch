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
	"strings"
	"sync"
	"time"

	"github.com/circonus-labs/circonus-logwatch/internal/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

// Statsd defines the relevant properties of a StatsD connection.
type Statsd struct {
	logger zerolog.Logger
	conn   net.Conn
	id     string
	port   string
	prefix string
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

// New initializes a udp connection to the localhost.
func New() (*Statsd, error) {
	id := viper.GetString(config.KeyDestCfgID)
	if id == "" {
		return nil, fmt.Errorf("invalid id, empty")
	}

	port := viper.GetString(config.KeyDestCfgPort)
	if port == "" {
		return nil, fmt.Errorf("invalid port, empty")
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

// Start the statsd Statsd.
func (c *Statsd) Start() error {
	return c.open()
}

// Stop the statsd Statsd.
func (c *Statsd) Stop() error {
	if c.conn == nil {
		return nil
	}
	return c.conn.Close()
}

// SetGaugeValue sends a gauge metric.
func (c *Statsd) SetGaugeValue(metric string, value interface{}) error { // gauge (ints or floats)
	v, err := getGaugeValue(value)
	if err != nil {
		return err
	}
	return c.send(fmt.Sprintf("%s:%s|g", metric, v))
}

// SetGaugeValueWithTags sends a gauge metric.
func (c *Statsd) SetGaugeValueWithTags(metric string, tags []string, value interface{}) error { // gauge (ints or floats)
	v, err := getGaugeValue(value)
	if err != nil {
		return err
	}
	return c.send(fmt.Sprintf("%s:%s|g|#%s", metric, v, strings.Join(tags, ",")))
}

// SetTimingValue sends a timing metric.
func (c *Statsd) SetTimingValue(metric string, value float64) error { // histogram
	return c.SetHistogramValue(metric, value)
}

// SetTimingValueWithTags sends a timing metric.
func (c *Statsd) SetTimingValueWithTags(metric string, tags []string, value float64) error { // histogram
	return c.SetHistogramValueWithTags(metric, tags, value)
}

// SetHistogramValue sends a histogram metric.
func (c *Statsd) SetHistogramValue(metric string, value float64) error { // histogram
	return c.send(fmt.Sprintf("%s:%e|ms", metric, value))
}

// SetHistogramValueWithTags sends a histogram metric.
func (c *Statsd) SetHistogramValueWithTags(metric string, tags []string, value float64) error { // histogram
	return c.send(fmt.Sprintf("%s:%e|ms|#%s", metric, value, strings.Join(tags, ",")))
}

// IncrementCounter sends a counter increment.
func (c *Statsd) IncrementCounter(metric string) error { // counter (monotonically increasing value)
	return c.IncrementCounterByValue(metric, 1)
}

// IncrementCounterWithTags sends a counter increment.
func (c *Statsd) IncrementCounterWithTags(metric string, tags []string) error { // counter (monotonically increasing value)
	return c.IncrementCounterByValueWithTags(metric, tags, 1)
}

// IncrementCounterByValue sends value to add to counter.
func (c *Statsd) IncrementCounterByValue(metric string, value uint64) error { // counter (monotonically increasing value)
	return c.send(fmt.Sprintf("%s:%d|c", metric, value))
}

// IncrementCounterByValueWithTags sends value to add to counter.
func (c *Statsd) IncrementCounterByValueWithTags(metric string, tags []string, value uint64) error { // counter (monotonically increasing value)
	return c.send(fmt.Sprintf("%s:%d|c|#%s", metric, value, strings.Join(tags, ",")))
}

// AddSetValue sends a unique value to the set metric.
func (c *Statsd) AddSetValue(metric string, value string) error { // set metric (ala statsd, counts unique values)
	return c.send(fmt.Sprintf("%s:%s|s", metric, value))
}

// AddSetValueWithTags sends a unique value to the set metric.
func (c *Statsd) AddSetValueWithTags(metric string, tags []string, value string) error { // set metric (ala statsd, counts unique values)
	return c.send(fmt.Sprintf("%s:%s|s|#%s", metric, value, strings.Join(tags, ",")))
}

// SetTextValue sends a text metric.
func (c *Statsd) SetTextValue(metric string, value string) error { // text metric
	return c.send(fmt.Sprintf("%s:%s|t", metric, value))
}

// SetTextValueWithTags sends a text metric.
func (c *Statsd) SetTextValueWithTags(metric string, tags []string, value string) error { // text metric
	return c.send(fmt.Sprintf("%s:%s|t|#%s", metric, value, strings.Join(tags, ",")))
}

// send stats data to udp statsd daemon
//
// Outgoing metric format:
//
//	name:value|type[|#tags]
//
// e.g.
//
//	foo:1|c
//	foo:1|c|#foo:bar
//	bar:2.5|ms
//	bar:2.5|ms|#foo:bar,baz:qux
//	baz:25|g
//	qux:abcd123|s
//	dib:38.282|h
//	dab:yadda yadda yadda|t
func (c *Statsd) send(metric string) error {
	if c.conn == nil {
		if err := c.open(); err != nil {
			return err
		}
	}

	m := c.prefix + metric

	c.logger.Debug().Str("metric", m).Msg("sending")

	_, err := fmt.Fprint(c.conn, m)
	if err != nil {
		return err
	}

	return nil
}

// open udp connection.
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

// getGaugeValue as string from interface.
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
		return "", fmt.Errorf("unknown type for value %v", v)
	}
	return vs, nil
}
