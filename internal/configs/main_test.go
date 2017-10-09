// Copyright © 2017 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package configs

import (
	"testing"

	"github.com/circonus-labs/circonus-logwatch/internal/config"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

func TestLoadConfigs(t *testing.T) {
	t.Log("Testing LoadConfigs")

	zerolog.SetGlobalLevel(zerolog.Disabled)

	t.Log("no config")
	{
		_, err := Load()
		if err == nil {
			t.Fatal("expected error")
		}
	}

	t.Log("no entries")
	{
		viper.Set(config.KeyLogConfDir, "testdata/empty")
		_, err := Load()
		viper.Reset()

		if err == nil {
			t.Fatal("expected error")
		}
	}

	t.Log("valid")
	{
		viper.Set(config.KeyLogConfDir, "testdata/")
		cfgs, err := Load()
		viper.Reset()

		if err != nil {
			t.Fatalf("expected no error, got (%s)", err)
		}

		if len(cfgs) == 0 {
			t.Fatal("expected >0 configs")
		}

		t.Logf("%#v\n", cfgs[0])
		for i, m := range cfgs[0].Metrics {
			t.Logf("\trule: %d = %#v\n", i, m)
		}
	}
}
