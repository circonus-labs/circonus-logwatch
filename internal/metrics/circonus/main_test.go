// Copyright © 2017 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package circonus

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/circonus-labs/circonus-logwatch/internal/config"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

var ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, client")
}))

func TestNew(t *testing.T) {
	t.Log("Testing New")

	t.Log("invalid")
	{
		_, err := New()
		if err == nil {
			t.Fatal("expected error")
		}
	}

	t.Log("invalid, agent (no url)")
	{
		viper.Set(config.KeyDestType, "agent")
		_, err := New()
		if err == nil {
			t.Fatal("expected error")
		}
		viper.Reset()
	}
	t.Log("invalid, agent (invalid interval)")
	{
		viper.Set(config.KeyDestType, "agent")
		viper.Set(config.KeyDestCfgAgentInterval, "-1")
		_, err := New()
		if err == nil {
			t.Fatal("expected error")
		}
		viper.Reset()
	}

	t.Log("valid, agent")
	{
		viper.Set(config.KeyDestType, "agent")
		viper.Set(config.KeyDestAgentURL, ts.URL)
		viper.Set(config.KeyDestCfgAgentInterval, "10s")
		_, err := New()
		if err != nil {
			t.Fatalf("expected no error, got (%s)", err)
		}
		viper.Reset()
	}

	t.Log("valid, check")
	{
		viper.Set(config.KeyDestType, "check")
		viper.Set(config.KeyDestCfgURL, ts.URL)
		_, err := New()
		if err != nil {
			t.Fatalf("expected no error, got (%s)", err)
		}
		viper.Reset()
	}
}

func TestStart(t *testing.T) {
	t.Log("Testing Start")

	t.Log("agent")
	{
		viper.Set(config.KeyDestType, "agent")
		viper.Set(config.KeyDestAgentURL, ts.URL)
		c, err := New()
		if err != nil {
			t.Fatalf("expected no error, got (%s)", err)
		}
		if err := c.Start(); err != nil {
			t.Fatalf("expected no error, got (%s)", err)
		}
		viper.Reset()
	}

	t.Log("check")
	{
		viper.Set(config.KeyDestType, "check")
		viper.Set(config.KeyDestCfgURL, ts.URL)
		c, err := New()
		if err != nil {
			t.Fatalf("expected no error, got (%s)", err)
		}
		if err := c.Start(); err != nil {
			t.Fatalf("expected no error, got (%s)", err)
		}
		viper.Reset()
	}
}

func TestStop(t *testing.T) {
	t.Log("Testing Stop")

	t.Log("agent")
	{
		viper.Set(config.KeyDestType, "agent")
		viper.Set(config.KeyDestAgentURL, ts.URL)
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

	t.Log("check")
	{
		viper.Set(config.KeyDestType, "check")
		viper.Set(config.KeyDestCfgURL, ts.URL)
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
}

func TestIncrementCounter(t *testing.T) {
	t.Log("Testing IncrementCounter")
	zerolog.SetGlobalLevel(zerolog.Disabled)

	t.Log("agent")
	{
		viper.Set(config.KeyDestType, "agent")
		viper.Set(config.KeyDestAgentURL, ts.URL)
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

	t.Log("check")
	{
		viper.Set(config.KeyDestType, "check")
		viper.Set(config.KeyDestCfgURL, ts.URL)
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
}

func TestIncrementCounterByValue(t *testing.T) {
	t.Log("Testing IncrementCounterByValue")
	zerolog.SetGlobalLevel(zerolog.Disabled)

	t.Log("agent")
	{
		viper.Set(config.KeyDestType, "agent")
		viper.Set(config.KeyDestAgentURL, ts.URL)
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

	t.Log("check")
	{
		viper.Set(config.KeyDestType, "check")
		viper.Set(config.KeyDestCfgURL, ts.URL)
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
}

func TestSetGaugeValue(t *testing.T) {
	t.Log("Testing SetGaugeValue")
	zerolog.SetGlobalLevel(zerolog.Disabled)

	t.Log("agent")
	{
		viper.Set(config.KeyDestType, "agent")
		viper.Set(config.KeyDestAgentURL, ts.URL)
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

	t.Log("check")
	{
		viper.Set(config.KeyDestType, "check")
		viper.Set(config.KeyDestCfgURL, ts.URL)
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
}

func TestSetHistogramValue(t *testing.T) {
	t.Log("Testing SetHistogramValue")
	zerolog.SetGlobalLevel(zerolog.Disabled)

	t.Log("agent")
	{
		viper.Set(config.KeyDestType, "agent")
		viper.Set(config.KeyDestAgentURL, ts.URL)
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

	t.Log("check")
	{
		viper.Set(config.KeyDestType, "check")
		viper.Set(config.KeyDestCfgURL, ts.URL)
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
}

func TestSetTimingValue(t *testing.T) {
	t.Log("Testing SetTimingValue")
	zerolog.SetGlobalLevel(zerolog.Disabled)

	t.Log("agent")
	{
		viper.Set(config.KeyDestType, "agent")
		viper.Set(config.KeyDestAgentURL, ts.URL)
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

	t.Log("check")
	{
		viper.Set(config.KeyDestType, "check")
		viper.Set(config.KeyDestCfgURL, ts.URL)
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
}

func TestAddSetValue(t *testing.T) {
	t.Log("Testing AddSetValue")
	zerolog.SetGlobalLevel(zerolog.Disabled)

	t.Log("agent")
	{
		viper.Set(config.KeyDestType, "agent")
		viper.Set(config.KeyDestAgentURL, ts.URL)
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

	t.Log("check")
	{
		viper.Set(config.KeyDestType, "check")
		viper.Set(config.KeyDestCfgURL, ts.URL)
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
}

func TestSetTextValue(t *testing.T) {
	t.Log("Testing SetTextValue")
	zerolog.SetGlobalLevel(zerolog.Disabled)

	t.Log("agent")
	{
		viper.Set(config.KeyDestType, "agent")
		viper.Set(config.KeyDestAgentURL, ts.URL)
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

	t.Log("check")
	{
		viper.Set(config.KeyDestType, "check")
		viper.Set(config.KeyDestCfgURL, ts.URL)
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
}
