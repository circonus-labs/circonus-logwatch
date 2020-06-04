// Copyright © 2017 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package config

import (
	"encoding/json"
	"expvar"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/circonus-labs/circonus-logwatch/internal/config/defaults"
	"github.com/circonus-labs/circonus-logwatch/internal/release"
	toml "github.com/pelletier/go-toml"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	yaml "gopkg.in/yaml.v2"
)

// Log defines the running config.log structure
type Log struct {
	Level  string `json:"level" yaml:"level" toml:"level"`
	Pretty bool   `json:"pretty" yaml:"pretty" toml:"pretty"`
}

// API defines the running config.api structure
type API struct {
	Key    string `json:"key" yaml:"key" toml:"key"`
	App    string `json:"app" yaml:"app" toml:"app"`
	URL    string `json:"url" yaml:"url" toml:"url"`
	CAFile string `mapstructure:"ca_file" json:"ca_file" yaml:"ca_file" toml:"ca_file"`
}

// DestConfig defines the running config.destination.config structure
type DestConfig struct {
	ID            string `json:"id" yaml:"id" toml:"id"`
	Port          string `json:"port" yaml:"port" toml:"port"`
	StatsdPrefix  string `mapstructure:"statsd_prefix" json:"statsd_prefix" yaml:"statsd_prefix" toml:"statsd_prefix"`
	CID           string `json:"cid" yaml:"cid" toml:"cid"`
	URL           string `json:"url" yaml:"url" toml:"url"`
	Target        string `json:"target" yaml:"target" toml:"target"`
	SearchTag     string `mapstructure:"search_tag" json:"search_tag" yaml:"search_tag" toml:"search_tag"`
	InstanceID    string `mapstructure:"instance_id" json:"instance_id" yaml:"instance_id" toml:"instance_id"`
	AgentInterval string `mapstructure:"agent_interval" json:"agent_interval" toml:"agent_interval" yaml:"agent_interval"`
}

// Destination defines the running config.destination structure
type Destination struct {
	Type   string     `json:"type" yaml:"type" toml:"type"`
	Config DestConfig `json:"config" yaml:"config" toml:"config"`
}

// Config defines the running config structure
type Config struct {
	LogConfDir  string      `mapstructure:"log_conf_dir" json:"log_conf_dir" yaml:"log_conf_dir" toml:"log_conf_dir"`
	AppStatPort string      `mapstructure:"app_stat_port" json:"app_stat_port" yaml:"app_stat_port" toml:"app_stat_port"`
	Debug       bool        `json:"debug" yaml:"debug" toml:"debug"`
	DebugCGM    bool        `mapstructure:"debug_cgm" json:"debug_cgm" yaml:"debug_cgm" toml:"debug_cgm"`
	DebugTail   bool        `mapstructure:"debug_tail" json:"debug_tail" yaml:"debug_tail" toml:"debug_tail"`
	DebugMetric bool        `mapstructure:"debug_metric" json:"debug_metric" yaml:"debug_metric" toml:"debug_metric"`
	API         API         `json:"api" yaml:"api" toml:"api"`
	Destination Destination `json:"destination" yaml:"destination" toml:"destination"`
	Log         Log         `json:"log" yaml:"log" toml:"log"`
}

//
// NOTE: adding a Key* MUST be reflected in the Config structures above
//
const (
	// KeyAPICAFile custom ca for circonus api (e.g. inside)
	KeyAPICAFile = "api.ca_file"

	// KeyAPITokenApp circonus api token key application name
	KeyAPITokenApp = "api.app"

	// KeyAPITokenKey circonus api token key
	KeyAPITokenKey = "api.key"

	// KeyAPIURL custom circonus api url (e.g. inside)
	KeyAPIURL = "api.url"

	// KeyDebug enables debug messages
	KeyDebug = "debug"

	// KeyDebugCGM enables debug messages for circonus-gometrics
	KeyDebugCGM = "debug_cgm"

	// KeyDebugTail enables debug messages for log tailing
	KeyDebugTail = "debug_tail"

	// KeyDebugMetric enables detailed log line/metric rule debugging messages
	KeyDebugMetric = "debug_metric"

	// KeyAppStatPort on which to expose runtime stats (expvar)
	KeyAppStatPort = "app_stat_port"

	// KeyLogConfDir log configuration directory
	KeyLogConfDir = "log_conf_dir"

	// KeyLogLevel logging level (panic, fatal, error, warn, info, debug, disabled)
	KeyLogLevel = "log.level"

	// KeyLogPretty output formatted log lines (for running in foreground)
	KeyLogPretty = "log.pretty"

	// KeyShowConfig - show configuration and exit
	KeyShowConfig = "show-config"

	// KeyShowVersion - show version information and exit
	KeyShowVersion = "version"

	// KeyDestType of destination where metrics are being sent (none|statsd|agent|check)
	KeyDestType = "destination.type"

	// KeyDestCfgID for destination type (statsd|agent)
	KeyDestCfgID = "destination.config.id"

	// KeyDestCfgCID for destination type (check, check bundle id)
	KeyDestCfgCID = "destination.config.cid"

	// KeyDestCfgURL for destination type (check, submission url)
	KeyDestCfgURL = "destination.config.url"

	// KeyDestCfgPort for destination type (statsd|agent, port to use agent=2609, statsd=8125)
	KeyDestCfgPort = "destination.config.port"

	// KeyDestCfgTarget for destination type (check, target, for search|create)
	KeyDestCfgTarget = "destination.config.target"

	// KeyDestCfgSearchTag for destination type (check, tag, for search|create)
	KeyDestCfgSearchTag = "destination.config.search_tag"

	// KeyDestCfgInstanceID for destination type (check, instance id, for search|create)
	KeyDestCfgInstanceID = "destination.config.instance_id"

	// KeyDestCfgStatsdPrefix to prepend on every metric
	KeyDestCfgStatsdPrefix = "destination.config.statsd_prefix"

	// KeyDestCfgAgentInterval send metrics this often, parsed as a time.Duration (agent)
	KeyDestCfgAgentInterval = "destination.config.agent_interval"

	// KeyDestAgentURL defines the submission url for the agent destination
	// NOTE: this is dynamically created by config validation, it is NOT part of Config
	KeyDestAgentURL = "destination.agentURL"

	cosiName = "cosi"
)

var (
	cosiCfgFile = filepath.Join(defaults.BasePath, "..", cosiName, "etc", "cosi.json")
)

// Validate the configuration options supplied
func Validate() error {

	if err := logConfDir(); err != nil {
		return err
	}

	if err := destConf(); err != nil {
		return err
	}

	if viper.GetString(KeyDestType) == "check" {
		if err := apiConf(); err != nil {
			return err
		}
	}

	return nil
}

func destConf() error {
	dest := viper.GetString(KeyDestType)
	switch dest {
	case "log":
		return nil // nothing to validate

	case "check":
		return nil // cgm will vet the config

	case "statsd":
		id := viper.GetString(KeyDestCfgID)
		if id == "" {
			viper.Set(KeyDestCfgID, release.NAME)
		}
		port := viper.GetString(KeyDestCfgPort)
		if port == "" {
			port = defaults.StatsdPort
			viper.Set(KeyDestCfgPort, port)
		}

		addr := net.JoinHostPort("localhost", port)
		a, err := net.ResolveUDPAddr("udp", addr)
		if err != nil {
			return errors.Wrapf(err, "destination %s, port %s", dest, addr)
		}

		if err := testPort("udp", a.String()); err != nil {
			return errors.Wrapf(err, "destination %s, port %s", dest, addr)
		}

	case "agent":
		id := viper.GetString(KeyDestCfgID)
		if id == "" {
			viper.Set(KeyDestCfgID, release.NAME)
		}
		port := viper.GetString(KeyDestCfgPort)
		if port == "" {
			port = defaults.AgentPort
			viper.Set(KeyDestCfgPort, port)
		}

		addr := net.JoinHostPort("localhost", port)
		a, err := net.ResolveTCPAddr("tcp", addr)
		if err != nil {
			return errors.Wrapf(err, "destination %s, port %s", dest, addr)
		}

		if err := testPort("tcp", a.String()); err != nil {
			return errors.Wrapf(err, "destination %s, port %s", dest, addr)
		}

		viper.Set(KeyDestAgentURL, fmt.Sprintf("http://%s/write/%s", a.String(), id))

	default:
		return errors.Errorf("invalid/unknown metric destination (%s)", dest)
	}

	return nil
}

func apiConf() error {
	apiKey := viper.GetString(KeyAPITokenKey)
	apiApp := viper.GetString(KeyAPITokenApp)
	apiURL := viper.GetString(KeyAPIURL)

	// if key is 'cosi' - load the cosi api config
	if strings.ToLower(apiKey) == cosiName {
		cKey, cApp, cURL, err := loadCOSIConfig()
		if err != nil {
			return err
		}

		apiKey = cKey
		apiApp = cApp
		apiURL = cURL
	}

	// API is required for reverse and/or statsd

	if apiKey == "" {
		return errors.New("API key is required")
	}

	if apiApp == "" {
		return errors.New("API app is required")
	}

	if apiURL == "" {
		return errors.New("API URL is required")
	}

	if apiURL != defaults.APIURL {
		parsedURL, err := url.Parse(apiURL)
		if err != nil {
			return errors.Wrap(err, "Invalid API URL")
		}
		if parsedURL.Scheme == "" || parsedURL.Host == "" || parsedURL.Path == "" {
			return errors.Errorf("Invalid API URL (%s)", apiURL)
		}
	}

	viper.Set(KeyAPITokenKey, apiKey)
	viper.Set(KeyAPITokenApp, apiApp)
	viper.Set(KeyAPIURL, apiURL)

	return nil
}

type cosiConfig struct {
	APIKey string `json:"api_key"`
	APIApp string `json:"api_app"`
	APIURL string `json:"api_url"`
}

func loadCOSIConfig() (string, string, string, error) {
	data, err := ioutil.ReadFile(cosiCfgFile)
	if err != nil {
		return "", "", "", errors.Wrap(err, "Unable to access cosi config")
	}

	var cfg cosiConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return "", "", "", errors.Wrapf(err, "Unable to parse cosi config (%s)", cosiCfgFile)
	}

	if cfg.APIKey == "" {
		return "", "", "", errors.Errorf("Missing API key, invalid cosi config (%s)", cosiCfgFile)
	}
	if cfg.APIApp == "" {
		return "", "", "", errors.Errorf("Missing API app, invalid cosi config (%s)", cosiCfgFile)
	}
	if cfg.APIURL == "" {
		return "", "", "", errors.Errorf("Missing API URL, invalid cosi config (%s)", cosiCfgFile)
	}

	return cfg.APIKey, cfg.APIApp, cfg.APIURL, nil

}

func logConfDir() error {
	errMsg := "Invalid log configuration directory"
	dir := viper.GetString(KeyLogConfDir)

	if dir == "" {
		return errors.Errorf(errMsg+" (%s)", dir)
	}

	absDir, err := filepath.Abs(dir)
	if err != nil {
		return errors.Wrap(err, errMsg)
	}

	dir = absDir

	fi, err := os.Stat(dir)
	if err != nil {
		return errors.Wrap(err, errMsg)
	}

	if !fi.Mode().IsDir() {
		return errors.Errorf(errMsg+" (%s) not a directory", dir)
	}

	// also try opening, to verify permissions
	// if last dir on path is not accessible to user, stat doesn't return EPERM
	f, err := os.Open(dir)
	if err != nil {
		return errors.Wrap(err, errMsg)
	}
	f.Close()

	viper.Set(KeyLogConfDir, dir)

	return nil
}

// testPort is used to verify agent|statsd port
func testPort(network, address string) error {
	c, err := net.Dial(network, address)
	if err != nil {
		return err
	}

	return c.Close()
}

// StatConfig adds the running config to the app stats
func StatConfig() error {
	cfg, err := getConfig()
	if err != nil {
		return err
	}

	cfg.API.Key = "..."
	cfg.API.App = "..."

	expvar.Publish("config", expvar.Func(func() interface{} {
		return &cfg
	}))

	return nil
}

// getConfig dumps the current configuration and returns it
func getConfig() (*Config, error) {
	var cfg *Config

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, errors.Wrap(err, "parsing config")
	}

	return cfg, nil
}

// ShowConfig prints the running configuration
func ShowConfig(w io.Writer) error {
	var cfg *Config
	var err error
	var data []byte

	cfg, err = getConfig()
	if err != nil {
		return err
	}

	format := viper.GetString(KeyShowConfig)

	log.Debug().Str("format", format).Msg("show-config")

	switch format {
	case "json":
		data, err = json.MarshalIndent(cfg, " ", "  ")
		if err != nil {
			return errors.Wrap(err, "formatting config (json)")
		}
	case "yaml":
		data, err = yaml.Marshal(cfg)
		if err != nil {
			return errors.Wrap(err, "formatting config (yaml)")
		}
	case "toml":
		data, err = toml.Marshal(*cfg)
		if err != nil {
			return errors.Wrap(err, "formatting config (toml)")
		}
	default:
		return errors.Errorf("unknown config format '%s'", format)
	}

	fmt.Fprintf(w, "%s v%s running config:\n%s\n", release.NAME, release.VERSION, data)
	return nil
}
