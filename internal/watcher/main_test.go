// Copyright © 2017 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package watcher

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/circonus-labs/circonus-logwatch/internal/config"
	"github.com/circonus-labs/circonus-logwatch/internal/configs"
	"github.com/circonus-labs/circonus-logwatch/internal/metrics/logonly"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

func fakeLog() {
	f, err := os.OpenFile("testdata/test.log", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		panic(err)
	}
	lines := []string{
		"testcounter\n",
		"testcounterval 2\n",
		"gaugefloat 1.23\n",
		"gaugeint 22\n",
		"hist 3.86\n",
		"timing 124.9\n",
		"set foo\n",
		"text|foo bar baz\n",
		"bad_type\n",
	}
	for {
		for _, line := range lines {
			if _, err := f.WriteString(line); err != nil {
				fmt.Printf("error writing to file %s", err)
				break
			}
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func TestNew(t *testing.T) {
	t.Log("Testing New")

	zerolog.SetGlobalLevel(zerolog.Disabled)

	t.Log("no dest, no log conf")
	{
		_, err := New(nil, nil)
		if err == nil {
			t.Fatal("expected error")
		}
	}

	t.Log("dest, no log conf")
	{
		viper.Set(config.KeyDestType, "log")
		dest, err := logonly.New()
		if err != nil {
			t.Fatalf("expected no error, got (%s)", err)
		}
		if _, err := New(dest, nil); err == nil {
			t.Fatal("expected error")
		}
		viper.Reset()
	}

	t.Log("no dest, log conf")
	{
		lc := &configs.Config{}
		if _, err := New(nil, lc); err == nil {
			t.Fatal("expected error")
		}
	}

	t.Log("valid")
	{
		viper.Set(config.KeyDestType, "log")
		lc := &configs.Config{ID: "test", LogFile: "testdata/test.log"}
		dest, err := logonly.New()
		if err != nil {
			t.Fatalf("expected no error, got (%s)", err)
		}
		if _, err := New(dest, lc); err != nil {
			t.Fatalf("expected no error, got (%s)", err)
		}
		viper.Reset()
	}
}

func TestStart(t *testing.T) {
	t.Log("Testing Start")

	zerolog.SetGlobalLevel(zerolog.Disabled)

	viper.Set(config.KeyDestType, "log")
	lc := &configs.Config{ID: "test", LogFile: "testdata/test.log"}
	dest, err := logonly.New()
	if err != nil {
		t.Fatalf("expected no error, got (%s)", err)
	}
	w, err := New(dest, lc)
	if err != nil {
		t.Fatalf("expected no error, got (%s)", err)
	}
	time.AfterFunc(1*time.Second, func() {
		w.Stop()
	})
	if err := w.Start(); err != nil {
		t.Fatalf("expected no error, got (%s)", err)
	}
	viper.Reset()
}

func TestStop(t *testing.T) {
	t.Log("Testing Stop")

	zerolog.SetGlobalLevel(zerolog.Disabled)

	viper.Set(config.KeyDestType, "log")
	lc := &configs.Config{ID: "test", LogFile: "testdata/test.log"}
	dest, err := logonly.New()
	if err != nil {
		t.Fatalf("expected no error, got (%s)", err)
	}
	w, err := New(dest, lc)
	if err != nil {
		t.Fatalf("expected no error, got (%s)", err)
	}
	time.AfterFunc(1*time.Second, func() {
		w.Stop()
	})
	if err := w.Start(); err != nil {
		t.Fatalf("expected no error, got (%s)", err)
	}
	viper.Reset()
}

func TestFull(t *testing.T) {
	t.Log("Testing full cycle")

	zerolog.SetGlobalLevel(zerolog.Disabled)

	go fakeLog()
	viper.Set(config.KeyDestType, "log")
	viper.Set(config.KeyLogConfDir, "testdata")
	cfgs, err := configs.Load()
	if err != nil {
		t.Fatalf("expected no error, got %s", err)
	}
	if len(cfgs) != 1 {
		t.Fatalf("expected 1 log config, got %d", len(cfgs))
	}
	dest, err := logonly.New()
	if err != nil {
		t.Fatalf("expected no error, got (%s)", err)
	}
	w, err := New(dest, cfgs[0])
	if err != nil {
		t.Fatalf("expected no error, got (%s)", err)
	}
	time.AfterFunc(5*time.Second, func() {
		w.Stop()
	})
	if err := w.Start(); err != nil {
		t.Fatalf("expected no error, got (%s)", err)
	}
	viper.Reset()
}
