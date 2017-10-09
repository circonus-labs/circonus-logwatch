// Copyright © 2017 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package statsd

import (
	"testing"

	"github.com/circonus-labs/circonus-logwatch/internal/config"
	"github.com/circonus-labs/circonus-logwatch/internal/config/defaults"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

func TestNew(t *testing.T) {
	t.Log("Testing New")

	t.Log("no config")
	{
		_, err := New()
		if err == nil {
			t.Fatal("expected error")
		}
	}

	t.Log("id, no port")
	{
		viper.Set(config.KeyDestCfgID, "foo")
		_, err := New()
		if err == nil {
			t.Fatal("expected error")
		}
		viper.Reset()
	}

	t.Log("valid")
	{
		viper.Set(config.KeyDestCfgID, "foo")
		viper.Set(config.KeyDestCfgPort, defaults.StatsdPort)
		_, err := New()
		if err != nil {
			t.Fatalf("expected no error, got (%s)", err)
		}
		viper.Reset()
	}
}

func TestStart(t *testing.T) {
	t.Log("Testing Start")

	viper.Set(config.KeyDestCfgID, "foo")
	viper.Set(config.KeyDestCfgPort, defaults.StatsdPort)

	c, err := New()
	if err != nil {
		t.Fatalf("expected no error, got (%s)", err)
	}

	if err := c.Start(); err != nil {
		t.Fatalf("expected no error, got (%s)", err)
	}

	viper.Reset()
}

func TestStop(t *testing.T) {
	t.Log("Testing Stop")

	viper.Set(config.KeyDestCfgID, "foo")
	viper.Set(config.KeyDestCfgPort, defaults.StatsdPort)

	c, err := New()
	if err != nil {
		t.Fatalf("expected no error, got (%s)", err)
	}

	if err := c.Start(); err != nil {
		t.Fatalf("expected no error, got (%s)", err)
	}

	if err := c.Stop(); err != nil {
		t.Fatalf("expected no error, got (%s)", err)
	}

	viper.Reset()
}

func TestIncrementCounter(t *testing.T) {
	t.Log("Testing IncrementCounter")
	zerolog.SetGlobalLevel(zerolog.Disabled)

	viper.Set(config.KeyDestCfgID, "foo")
	viper.Set(config.KeyDestCfgPort, defaults.StatsdPort)

	c, err := New()
	if err != nil {
		t.Fatalf("expected no error, got (%s)", err)
	}

	if err := c.Start(); err != nil {
		t.Fatalf("expected no error, got (%s)", err)
	}

	if err := c.IncrementCounter("foo"); err != nil {
		t.Fatalf("expected no error, got (%s)", err)
	}

	if err := c.Stop(); err != nil {
		t.Fatalf("expected no error, got (%s)", err)
	}

	viper.Reset()
}

func TestIncrementCounterByValue(t *testing.T) {
	t.Log("Testing IncrementCounterByValue")
	zerolog.SetGlobalLevel(zerolog.Disabled)

	viper.Set(config.KeyDestCfgID, "foo")
	viper.Set(config.KeyDestCfgPort, defaults.StatsdPort)

	c, err := New()
	if err != nil {
		t.Fatalf("expected no error, got (%s)", err)
	}

	if err := c.Start(); err != nil {
		t.Fatalf("expected no error, got (%s)", err)
	}

	if err := c.IncrementCounterByValue("foo", 1); err != nil {
		t.Fatalf("expected no error, got (%s)", err)
	}

	if err := c.Stop(); err != nil {
		t.Fatalf("expected no error, got (%s)", err)
	}

	viper.Reset()
}

func TestSetGaugeValue(t *testing.T) {
	t.Log("Testing SetGaugeValue")
	zerolog.SetGlobalLevel(zerolog.Disabled)

	viper.Set(config.KeyDestCfgID, "foo")
	viper.Set(config.KeyDestCfgPort, defaults.StatsdPort)

	c, err := New()
	if err != nil {
		t.Fatalf("expected no error, got (%s)", err)
	}

	if err := c.Start(); err != nil {
		t.Fatalf("expected no error, got (%s)", err)
	}

	if err := c.SetGaugeValue("foo", 1); err != nil {
		t.Fatalf("expected no error, got (%s)", err)
	}

	if err := c.Stop(); err != nil {
		t.Fatalf("expected no error, got (%s)", err)
	}

	viper.Reset()
}

func TestSetHistogramValue(t *testing.T) {
	t.Log("Testing SetHistogramValue")
	zerolog.SetGlobalLevel(zerolog.Disabled)

	viper.Set(config.KeyDestCfgID, "foo")
	viper.Set(config.KeyDestCfgPort, defaults.StatsdPort)

	c, err := New()
	if err != nil {
		t.Fatalf("expected no error, got (%s)", err)
	}

	if err := c.Start(); err != nil {
		t.Fatalf("expected no error, got (%s)", err)
	}

	if err := c.SetHistogramValue("foo", 1.23); err != nil {
		t.Fatalf("expected no error, got (%s)", err)
	}

	if err := c.Stop(); err != nil {
		t.Fatalf("expected no error, got (%s)", err)
	}

	viper.Reset()
}

func TestSetTimingValue(t *testing.T) {
	t.Log("Testing SetTimingValue")
	zerolog.SetGlobalLevel(zerolog.Disabled)

	viper.Set(config.KeyDestCfgID, "foo")
	viper.Set(config.KeyDestCfgPort, defaults.StatsdPort)

	c, err := New()
	if err != nil {
		t.Fatalf("expected no error, got (%s)", err)
	}

	if err := c.Start(); err != nil {
		t.Fatalf("expected no error, got (%s)", err)
	}

	if err := c.SetTimingValue("foo", 1.23); err != nil {
		t.Fatalf("expected no error, got (%s)", err)
	}

	if err := c.Stop(); err != nil {
		t.Fatalf("expected no error, got (%s)", err)
	}

	viper.Reset()
}

func TestAddSetValue(t *testing.T) {
	t.Log("Testing AddSetValue")
	zerolog.SetGlobalLevel(zerolog.Disabled)

	viper.Set(config.KeyDestCfgID, "foo")
	viper.Set(config.KeyDestCfgPort, defaults.StatsdPort)

	c, err := New()
	if err != nil {
		t.Fatalf("expected no error, got (%s)", err)
	}

	if err := c.Start(); err != nil {
		t.Fatalf("expected no error, got (%s)", err)
	}

	if err := c.AddSetValue("foo", "bar"); err != nil {
		t.Fatalf("expected no error, got (%s)", err)
	}

	if err := c.Stop(); err != nil {
		t.Fatalf("expected no error, got (%s)", err)
	}

	viper.Reset()
}

func TestSetTextValue(t *testing.T) {
	t.Log("Testing SetTextValue")
	zerolog.SetGlobalLevel(zerolog.Disabled)

	viper.Set(config.KeyDestCfgID, "foo")
	viper.Set(config.KeyDestCfgPort, defaults.StatsdPort)

	c, err := New()
	if err != nil {
		t.Fatalf("expected no error, got (%s)", err)
	}

	if err := c.Start(); err != nil {
		t.Fatalf("expected no error, got (%s)", err)
	}

	if err := c.SetTextValue("foo", "bar"); err != nil {
		t.Fatalf("expected no error, got (%s)", err)
	}

	if err := c.Stop(); err != nil {
		t.Fatalf("expected no error, got (%s)", err)
	}

	viper.Reset()
}

func TestGetGaugeValue(t *testing.T) {
	t.Log("Testing getGaugeValue")

	t.Log("int")
	if _, err := getGaugeValue(int(1)); err != nil {
		t.Fatalf("expected no error, got (%s)", err)
	}

	t.Log("int8")
	if _, err := getGaugeValue(int8(1)); err != nil {
		t.Fatalf("expected no error, got (%s)", err)
	}

	t.Log("int16")
	if _, err := getGaugeValue(int16(1)); err != nil {
		t.Fatalf("expected no error, got (%s)", err)
	}

	t.Log("int32")
	if _, err := getGaugeValue(int32(1)); err != nil {
		t.Fatalf("expected no error, got (%s)", err)
	}

	t.Log("int64")
	if _, err := getGaugeValue(int64(1)); err != nil {
		t.Fatalf("expected no error, got (%s)", err)
	}

	t.Log("uint")
	if _, err := getGaugeValue(uint(1)); err != nil {
		t.Fatalf("expected no error, got (%s)", err)
	}

	t.Log("uint8")
	if _, err := getGaugeValue(uint8(1)); err != nil {
		t.Fatalf("expected no error, got (%s)", err)
	}

	t.Log("uint16")
	if _, err := getGaugeValue(uint16(1)); err != nil {
		t.Fatalf("expected no error, got (%s)", err)
	}

	t.Log("uint32")
	if _, err := getGaugeValue(uint32(1)); err != nil {
		t.Fatalf("expected no error, got (%s)", err)
	}

	t.Log("uint64")
	if _, err := getGaugeValue(uint64(1)); err != nil {
		t.Fatalf("expected no error, got (%s)", err)
	}

	t.Log("float32")
	if _, err := getGaugeValue(float32(1.23)); err != nil {
		t.Fatalf("expected no error, got (%s)", err)
	}

	t.Log("float64")
	if _, err := getGaugeValue(float64(1.23)); err != nil {
		t.Fatalf("expected no error, got (%s)", err)
	}

	t.Log("invalid")
	if _, err := getGaugeValue(true); err == nil {
		t.Fatal("expected error")
	}
}
