// Copyright © 2017 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package config

import (
	"errors"
	"io/ioutil"
	"net"
	"path/filepath"
	"strings"
	"testing"

	"github.com/circonus-labs/circonus-logwatch/internal/config/defaults"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

func TestValidate(t *testing.T) {
	t.Log("Testing Validate")

	t.Log("no config")
	{
		err := Validate()
		if err == nil {
			t.Fatal("expected error")
		}
	}

	t.Log("log conf dir")
	{
		viper.Set(KeyLogConfDir, "testdata/")
		err := Validate()
		if err == nil {
			t.Fatal("expected error")
		}
		viper.Reset()
	}

	t.Log("log conf dir, dest log")
	{
		viper.Set(KeyDestType, "log")
		viper.Set(KeyLogConfDir, "testdata/")
		err := Validate()
		if err != nil {
			t.Fatalf("expected no error, got (%s)", err)
		}
		viper.Reset()
	}

	t.Log("log conf dir, dest statsd")
	{
		viper.Set(KeyDestType, "statsd")
		viper.Set(KeyLogConfDir, "testdata/")
		err := Validate()
		if err != nil {
			t.Fatalf("expected no error, got (%s)", err)
		}
		viper.Reset()
	}

	t.Log("log conf dir, dest agent")
	{
		ts, _ := net.Listen("tcp", "127.0.0.1:2609")
		viper.Set(KeyDestType, "agent")
		viper.Set(KeyLogConfDir, "testdata/")
		err := Validate()
		if err != nil {
			t.Fatalf("expected no error, got (%s)", err)
		}
		viper.Reset()
		if ts != nil {
			ts.Close()
		}
	}

	t.Log("log conf dir, dest check, missing api")
	{
		viper.Set(KeyDestType, "check")
		viper.Set(KeyLogConfDir, "testdata/")
		err := Validate()
		if err == nil {
			t.Fatal("expected error")
		}
		viper.Reset()
	}

	t.Log("log conf dir, dest check")
	{
		viper.Set(KeyDestType, "check")
		viper.Set(KeyAPITokenKey, "foo")
		viper.Set(KeyAPITokenApp, "foo")
		viper.Set(KeyAPIURL, defaults.APIURL)
		viper.Set(KeyLogConfDir, "testdata/")
		err := Validate()
		if err != nil {
			t.Fatalf("expected no error, got (%s)", err)
		}
		viper.Reset()
	}
}

func TestDestConf(t *testing.T) {
	t.Log("Testing destConf")

	t.Log("no config")
	{
		err := destConf()
		if err == nil {
			t.Fatal("expected error")
		}
	}

	t.Log("log")
	{
		viper.Set(KeyDestType, "log")
		err := destConf()
		if err != nil {
			t.Fatalf("expected no error, got (%s)", err)
		}
		viper.Reset()
	}

	t.Log("check")
	{
		viper.Set(KeyDestType, "check")
		err := destConf()
		if err != nil {
			t.Fatalf("expected no error, got (%s)", err)
		}
		viper.Reset()
	}

	t.Log("statsd, bad port")
	{
		viper.Set(KeyDestType, "statsd")
		viper.Set(KeyDestCfgPort, "foo")
		err := destConf()
		if err == nil {
			t.Fatal("expected error")
		}
		viper.Reset()
	}

	t.Log("statsd")
	{
		viper.Set(KeyDestType, "statsd")
		err := destConf()
		if err != nil {
			t.Fatalf("expected no error, got (%s)", err)
		}
		viper.Reset()
	}

	t.Log("agent, bad port")
	{
		viper.Set(KeyDestType, "agent")
		viper.Set(KeyDestCfgPort, "foo")
		if err := destConf(); err == nil {
			t.Fatal("expected error")
		}
		viper.Reset()
	}

	t.Log("agent, not listening")
	{
		viper.Set(KeyDestType, "agent")
		viper.Set(KeyDestCfgPort, "12345")
		if err := destConf(); err == nil {
			t.Fatal("expected error")
		}
		viper.Reset()
	}

	t.Log("agent, listening")
	{
		ts, _ := net.Listen("tcp", "127.0.0.1:2609")
		viper.Set(KeyDestType, "agent")
		if err := destConf(); err != nil {
			t.Fatalf("expected no error, got (%s)", err)
		}
		viper.Reset()
		if ts != nil {
			ts.Close()
		}
	}

}

func TestLogConfDir(t *testing.T) {
	t.Log("Testing logConfDir")
	zerolog.SetGlobalLevel(zerolog.Disabled)

	t.Log("no directory")
	{
		viper.Set(KeyLogConfDir, "")
		expectedError := errors.New("invalid log configuration directory ()")
		err := logConfDir()
		if err == nil {
			t.Fatalf("Expected error")
		}
		if err.Error() != expectedError.Error() {
			t.Errorf("Expected (%s) got (%s)", expectedError, err)
		}
	}

	t.Log("Invalid directory (not found)")
	{
		viper.Set(KeyLogConfDir, "foo")
		err := logConfDir()
		if err == nil {
			t.Fatalf("Expected error")
		}
		sfx := "internal/config/foo: no such file or directory"
		if !strings.HasSuffix(err.Error(), sfx) {
			t.Errorf("Expected (%s) got (%s)", sfx, err)
		}
	}

	// NOTE touch testdata/not_a_dir
	t.Log("Invalid directory (not a dir)")
	{
		viper.Set(KeyLogConfDir, filepath.Join("testdata", "not_a_dir"))
		err := logConfDir()
		if err == nil {
			t.Fatalf("Expected error")
		}
		sfx := "internal/config/testdata/not_a_dir) not a directory"
		if !strings.HasSuffix(err.Error(), sfx) {
			t.Errorf("Expected (%s) got (%s)", sfx, err)
		}
	}

	//
	// NOTE next two will fail if the directory structure isn't set up correctly (which is not 'git'able)
	//
	// sudo mkdir -p testdata/no_access_dir/test && sudo chmod -R 700 testdata/no_access_dir
	//
	// t.Log("Invalid directory (perms, subdir)")
	// {
	// 	viper.Set(KeyLogConfDir, filepath.Join("testdata", "no_access_dir", "test"))
	// 	err := logConfDir()
	// 	if err == nil {
	// 		t.Fatalf("Expected error - check 'sudo mkdir -p testdata/no_access_dir/test && sudo chmod -R 700 testdata/no_access_dir'")
	// 	}
	// 	sfx := "internal/config/testdata/no_access_dir/test: permission denied"
	// 	if !strings.HasSuffix(err.Error(), sfx) {
	// 		t.Errorf("Expected (%s) got (%s)", sfx, err)
	// 	}
	// }
	//
	// t.Log("Invalid directory (perms, open)")
	// {
	// 	viper.Set(KeyLogConfDir, filepath.Join("testdata", "no_access_dir"))
	// 	err := logConfDir()
	// 	if err == nil {
	// 		t.Fatalf("Expected error")
	// 	}
	// 	sfx := "internal/config/testdata/no_access_dir: permission denied"
	// 	if !strings.HasSuffix(err.Error(), sfx) {
	// 		t.Errorf("Expected (%s) got (%s)", sfx, err)
	// 	}
	// }

	t.Log("Valid directory")
	{
		viper.Set(KeyLogConfDir, "testdata")
		err := logConfDir()
		if err != nil {
			t.Fatal("Expected NO error")
		}
		dir := viper.GetString(KeyLogConfDir)
		sfx := "internal/config/testdata"
		if !strings.HasSuffix(dir, sfx) {
			t.Errorf("Expected (%s), got '%s'", sfx, dir)
		}
	}
}

func TestApiConf(t *testing.T) {
	t.Log("Testing apiConf")
	zerolog.SetGlobalLevel(zerolog.Disabled)

	t.Log("No key/app/url")
	{
		expectedError := errors.New("API key is required")
		err := apiConf()
		if err == nil {
			t.Fatal("Expected error")
		}
		if err.Error() != expectedError.Error() {
			t.Errorf("Expected (%s) got (%s)", expectedError, err)
		}
	}

	t.Log("key=cosi, no cfg")
	{
		viper.Set(KeyAPITokenKey, cosiName)
		err := apiConf()
		if err == nil {
			t.Fatal("Expected error")
		}
		pfx := "unable to access cosi config:"
		if !strings.HasPrefix(err.Error(), pfx) {
			t.Errorf("Expected (^%s) got (%s)", pfx, err)
		}
	}

	t.Log("No app")
	{
		viper.Set(KeyAPITokenKey, "foo")
		expectedError := errors.New("API app is required")
		err := apiConf()
		if err == nil {
			t.Fatal("Expected error")
		}
		if err.Error() != expectedError.Error() {
			t.Errorf("Expected (%s) got (%s)", expectedError, err)
		}
	}

	t.Log("No url")
	{
		viper.Set(KeyAPITokenKey, "foo")
		viper.Set(KeyAPITokenApp, "foo")
		expectedError := errors.New("API URL is required")
		err := apiConf()
		if err == nil {
			t.Fatal("Expected error")
		}
		if err.Error() != expectedError.Error() {
			t.Errorf("Expected (%s) got (%s)", expectedError, err)
		}
	}

	t.Log("Invalid url (foo)")
	{
		viper.Set(KeyAPITokenKey, "foo")
		viper.Set(KeyAPITokenApp, "foo")
		viper.Set(KeyAPIURL, "foo")
		expectedError := errors.New("invalid API URL (foo)")
		err := apiConf()
		if err == nil {
			t.Fatal("Expected error")
		}
		if err.Error() != expectedError.Error() {
			t.Errorf("Expected (%s) got (%s)", expectedError, err)
		}
	}

	t.Log("Invalid url (foo_bar://herp/derp)")
	{
		viper.Set(KeyAPITokenKey, "foo")
		viper.Set(KeyAPITokenApp, "foo")
		viper.Set(KeyAPIURL, "foo_bar://herp/derp")
		expectedError := errors.New(`invalid API URL: parse "foo_bar://herp/derp": first path segment in URL cannot contain colon`)
		err := apiConf()
		if err == nil {
			t.Fatal("Expected error")
		}
		if err.Error() != expectedError.Error() {
			t.Errorf("Expected (%s) got (%s)", expectedError, err)
		}
	}

	t.Log("Valid options")
	{
		viper.Set(KeyAPITokenKey, "foo")
		viper.Set(KeyAPITokenApp, "foo")
		viper.Set(KeyAPIURL, "http://foo.com/bar")
		err := apiConf()
		if err != nil {
			t.Fatalf("Expected NO error, got (%s)", err)
		}
	}

}

func TestShowConfig(t *testing.T) {
	t.Log("Testing ShowConfig")
	zerolog.SetGlobalLevel(zerolog.Disabled)

	t.Log("invalid")
	{
		viper.Set(KeyShowConfig, "invalid")
		err := ShowConfig(ioutil.Discard)
		if err == nil {
			t.Fatal("expected error")
		}
	}

	t.Log("YAML")
	{
		viper.Set(KeyShowConfig, "yaml")
		err := ShowConfig(ioutil.Discard)
		if err != nil {
			t.Fatalf("expected no error, got %s", err)
		}
	}

	t.Log("TOML")
	{
		viper.Set(KeyShowConfig, "toml")
		err := ShowConfig(ioutil.Discard)
		if err != nil {
			t.Fatalf("expected no error, got %s", err)
		}
	}

	t.Log("JSON")
	{
		viper.Set(KeyShowConfig, "json")
		err := ShowConfig(ioutil.Discard)
		if err != nil {
			t.Fatalf("expected no error, got %s", err)
		}
	}
}

func TestGetConfig(t *testing.T) {
	t.Log("Testing getConfig")
	zerolog.SetGlobalLevel(zerolog.Disabled)

	cfg, err := getConfig()
	if err != nil {
		t.Fatalf("expected no error, got %s", err)
	}
	if cfg == nil {
		t.Fatal("expected not nil")
	}
}

func TestStatConfig(t *testing.T) {
	t.Log("Testing StatConfig")
	zerolog.SetGlobalLevel(zerolog.Disabled)

	err := StatConfig()
	if err != nil {
		t.Fatalf("expected no error, got %s", err)
	}
}

func TestLoadCosiConfig(t *testing.T) {
	t.Log("testing loadCosiConfig")
	zerolog.SetGlobalLevel(zerolog.Disabled)
	t.Log("cosi - missing")
	{
		expectedError := errors.New("unable to access cosi config: open testdata/cosi_missing.json: no such file or directory")
		cosiCfgFile = filepath.Join("testdata", "cosi_missing.json")
		t.Logf("cosiCfgFile %s", cosiCfgFile)
		key, app, apiURL, err := loadCOSIConfig()
		if err == nil {
			t.Fatalf("Expected error")
		}
		if err.Error() != expectedError.Error() {
			t.Errorf("Expected (%s) got (%s)", expectedError, err)
		}
		if key != "" {
			t.Errorf("expected blank key")
		}
		if app != "" {
			t.Errorf("expected blank app")
		}
		if apiURL != "" {
			t.Errorf("expected blank url")
		}
	}

	t.Log("cosi - bad json")
	{
		expectedError := errors.New("Unable to parse cosi config (testdata/cosi_bad.json): invalid character '#' looking for beginning of value")
		cosiCfgFile = filepath.Join("testdata", "cosi_bad.json")
		t.Logf("cosiCfgFile %s", cosiCfgFile)
		key, app, apiURL, err := loadCOSIConfig()
		if err == nil {
			t.Fatalf("Expected error")
		}
		if err.Error() != expectedError.Error() {
			t.Errorf("Expected (%s) got (%s)", expectedError, err)
		}
		if key != "" {
			t.Errorf("expected blank key")
		}
		if app != "" {
			t.Errorf("expected blank app")
		}
		if apiURL != "" {
			t.Errorf("expected blank url")
		}
	}

	t.Log("cosi - invalid config missing key")
	{
		expectedError := errors.New("Missing API key, invalid cosi config (testdata/cosi_invalid_key.json)")
		cosiCfgFile = filepath.Join("testdata", "cosi_invalid_key.json")
		t.Logf("cosiCfgFile %s", cosiCfgFile)
		key, app, apiURL, err := loadCOSIConfig()
		if err == nil {
			t.Fatalf("Expected error")
		}
		if err.Error() != expectedError.Error() {
			t.Errorf("Expected (%s) got (%s)", expectedError, err)
		}
		if key != "" {
			t.Errorf("expected blank key")
		}
		if app != "" {
			t.Errorf("expected blank app")
		}
		if apiURL != "" {
			t.Errorf("expected blank url")
		}
	}

	t.Log("cosi - invalid config missing app")
	{
		expectedError := errors.New("Missing API app, invalid cosi config (testdata/cosi_invalid_app.json)")
		cosiCfgFile = filepath.Join("testdata", "cosi_invalid_app.json")
		t.Logf("cosiCfgFile %s", cosiCfgFile)
		key, app, apiURL, err := loadCOSIConfig()
		if err == nil {
			t.Fatalf("Expected error")
		}
		if err.Error() != expectedError.Error() {
			t.Errorf("Expected (%s) got (%s)", expectedError, err)
		}
		if key != "" {
			t.Errorf("expected blank key")
		}
		if app != "" {
			t.Errorf("expected blank app")
		}
		if apiURL != "" {
			t.Errorf("expected blank url")
		}
	}

	t.Log("cosi - invalid config missing url")
	{
		expectedError := errors.New("Missing API URL, invalid cosi config (testdata/cosi_invalid_url.json)")
		cosiCfgFile = filepath.Join("testdata", "cosi_invalid_url.json")
		t.Logf("cosiCfgFile %s", cosiCfgFile)
		key, app, apiURL, err := loadCOSIConfig()
		if err == nil {
			t.Fatalf("Expected error")
		}
		if err.Error() != expectedError.Error() {
			t.Errorf("Expected (%s) got (%s)", expectedError, err)
		}
		if key != "" {
			t.Errorf("expected blank key")
		}
		if app != "" {
			t.Errorf("expected blank app")
		}
		if apiURL != "" {
			t.Errorf("expected blank url")
		}
	}
}
