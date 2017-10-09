// Copyright © 2017 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package agent

import (
	"testing"

	"github.com/circonus-labs/circonus-logwatch/internal/config"
	"github.com/circonus-labs/circonus-logwatch/internal/config/defaults"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

func TestNew(t *testing.T) {
	t.Log("Testing New")

	zerolog.SetGlobalLevel(zerolog.Disabled)

	t.Log("no config")
	{
		_, err := New()
		if err == nil {
			t.Fatal("expected error")
		}
	}

	t.Log("no log configs")
	{
		viper.Set(config.KeyLogConfDir, "testdata/empty")
		viper.Set(config.KeyDestType, defaults.DestinationType)
		_, err := New()
		if err == nil {
			t.Fatal("expected error")
		}
		viper.Reset()
	}

	t.Log("valid w/defaults")
	{
		viper.Set(config.KeyLogConfDir, "testdata/")
		viper.Set(config.KeyDestType, defaults.DestinationType)
		a, err := New()
		if err != nil {
			t.Fatalf("expected no error, got %s", err)
		}
		if a == nil {
			t.Fatal("expected not nil")
		}
		viper.Reset()
	}
}

func TestStart(t *testing.T) {
	t.Skip("not testing Start")
}

func TestStop(t *testing.T) {
	t.Log("Testing Stop")

	zerolog.SetGlobalLevel(zerolog.Disabled)

	t.Log("valid w/defaults")
	{
		viper.Set(config.KeyLogConfDir, "testdata/")
		viper.Set(config.KeyDestType, defaults.DestinationType)
		a, err := New()
		if err != nil {
			t.Fatalf("expected no error, got %s", err)
		}
		if a == nil {
			t.Fatal("expected not nil")
		}

		a.Stop()
		viper.Reset()
	}
}
