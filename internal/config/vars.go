// Copyright © 2017 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package config

import (
	"path/filepath"

	"github.com/circonus-labs/circonus-logwatch/internal/config/defaults"
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
	ID           string `json:"id" yaml:"id" toml:"id"`
	Port         string `json:"port" yaml:"port" toml:"port"`
	StatsdPrefix string `mapstructure:"statsd_prefix" json:"statsd_prefix" yaml:"statsd_prefix" toml:"statsd_prefix"`
	CID          string `json:"cid" yaml:"cid" toml:"cid"`
	URL          string `json:"url" yaml:"url" toml:"url"`
	Target       string `json:"target" yaml:"target" toml:"target"`
	SearchTag    string `mapstructure:"search_tag" json:"search_tag" yaml:"search_tag" toml:"search_tag"`
	InstanceID   string `mapstructure:"instance_id" json:"instance_id" yaml:"instance_id" toml:"instance_id"`
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

	// KeyDestAgentURL defines the submission url for the agent destination
	// NOTE: this is dymanically created by config validation, it is NOT part of Config
	KeyDestAgentURL = "destination.agentURL"

	cosiName = "cosi"
)

var (
	cosiCfgFile = filepath.Join(defaults.BasePath, "..", cosiName, "etc", "cosi.json")
)
