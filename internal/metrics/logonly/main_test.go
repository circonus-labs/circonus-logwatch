// Copyright © 2017 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package logonly

import (
	"testing"

	"github.com/rs/zerolog"
)

func TestNew(t *testing.T) {
	t.Log("Testing New")

	_, err := New()
	if err != nil {
		t.Fatalf("expected no error, got (%s)", err)
	}
}

func TestStart(t *testing.T) {
	t.Log("Testing Start")

	c, err := New()
	if err != nil {
		t.Fatalf("expected no error, got (%s)", err)
	}

	if err := c.Start(); err != nil {
		t.Fatalf("expected no error, got (%s)", err)
	}
}

func TestStop(t *testing.T) {
	t.Log("Testing Stop")

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
}

func TestIncrementCounter(t *testing.T) {
	t.Log("Testing IncrementCounter")
	zerolog.SetGlobalLevel(zerolog.Disabled)

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
}

func TestIncrementCounterByValue(t *testing.T) {
	t.Log("Testing IncrementCounterByValue")
	zerolog.SetGlobalLevel(zerolog.Disabled)

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
}

func TestSetGaugeValue(t *testing.T) {
	t.Log("Testing SetGaugeValue")
	zerolog.SetGlobalLevel(zerolog.Disabled)

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
}

func TestSetHistogramValue(t *testing.T) {
	t.Log("Testing SetHistogramValue")
	zerolog.SetGlobalLevel(zerolog.Disabled)

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
}

func TestSetTimingValue(t *testing.T) {
	t.Log("Testing SetTimingValue")
	zerolog.SetGlobalLevel(zerolog.Disabled)

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
}

func TestAddSetValue(t *testing.T) {
	t.Log("Testing AddSetValue")
	zerolog.SetGlobalLevel(zerolog.Disabled)

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
}

func TestSetTextValue(t *testing.T) {
	t.Log("Testing SetTextValue")
	zerolog.SetGlobalLevel(zerolog.Disabled)

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
}
