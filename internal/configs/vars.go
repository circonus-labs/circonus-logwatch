// Copyright © 2017 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package configs

import (
	"regexp"
	"text/template"
)

// Metric is a metric definition for the log config
type Metric struct {
	Matcher    *regexp.Regexp
	MatchParts []string
	Namer      *template.Template
	ValueKey   string
	Match      string `json:"match" yaml:"match" toml:"match"`
	Name       string `json:"name" yaml:"name" toml:"name"`
	Type       string `json:"type" yaml:"type" toml:"type"`
}

// Config defines a log to watch
type Config struct {
	ID      string    `json:"id" yaml:"id" toml:"id"`
	LogFile string    `json:"log_file" yaml:"log_file" toml:"log_file"`
	Metrics []*Metric `json:"metrics" yaml:"metrics" toml:"metrics"`
}
