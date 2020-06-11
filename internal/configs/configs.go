// Copyright © 2017 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package configs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/circonus-labs/circonus-logwatch/internal/config"
	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

// Metric is a metric definition for the log config
type Metric struct {
	Matcher    *regexp.Regexp
	MatchParts []string
	Namer      *template.Template
	Tagger     *template.Template
	ValueKey   string
	Match      string `json:"match" yaml:"match" toml:"match"`
	Name       string `json:"name" yaml:"name" toml:"name"`
	Type       string `json:"type" yaml:"type" toml:"type"`
	Tags       string `json:"tags" toml:"tags" yaml:"tags"`
}

// Config defines a log to watch
type Config struct {
	ID      string    `json:"id" yaml:"id" toml:"id"`
	LogFile string    `json:"log_file" yaml:"log_file" toml:"log_file"`
	Metrics []*Metric `json:"metrics" yaml:"metrics" toml:"metrics"`
}

// Load reads the log configurations from log config directory
func Load() ([]*Config, error) {
	logger := log.With().Str("pkg", "configs").Logger()
	supportedConfExts := regexp.MustCompile(`^\.(yaml|json|toml)$`)
	logConfDir := viper.GetString(config.KeyLogConfDir)

	if logConfDir == "" {
		return nil, errors.Errorf("invalid log config directory (empty)")
	}

	logger.Debug().
		Str("dir", logConfDir).
		Msg("loading log configs")

	entries, err := ioutil.ReadDir(logConfDir)
	if err != nil {
		return nil, err
	}

	if len(entries) == 0 {
		return nil, errors.Errorf("no log configurations found in (%s)", logConfDir)
	}

	var cfgs []*Config

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		cfgFile := path.Join(logConfDir, entry.Name())
		cfgType := filepath.Ext(cfgFile)
		if !supportedConfExts.MatchString(cfgType) {
			logger.Warn().
				Str("type", cfgType).
				Str("file", cfgFile).
				Msg("unsupported config type, ignoring")
			continue
		}

		logger.Debug().
			Str("type", cfgType).
			Str("file", cfgFile).
			Msg("loading")
		logcfg, err := parse(cfgType, cfgFile)
		if err != nil {
			logger.Warn().
				Err(err).
				Str("file", cfgFile).
				Msg("parsing")
			continue
		}

		if err := checkLogFileAccess(logcfg.LogFile); err != nil {
			logger.Warn().
				Err(err).
				Str("log", logcfg.LogFile).
				Msg("access")
			continue
		}

		if logcfg.ID == "" { // ID not explicitly set, use the base of the config file name
			logcfg.ID = strings.Replace(filepath.Base(logcfg.LogFile), filepath.Ext(logcfg.LogFile), "", -1)
		}

		if validMetricRules(logcfg.ID, logger, logcfg.Metrics) {
			cfgs = append(cfgs, &logcfg)
		}
	}

	if len(cfgs) == 0 {
		return nil, errors.New("no valid configurations found")
	}

	return cfgs, nil
}

func validMetricRules(logID string, logger zerolog.Logger, rules []*Metric) bool {
	for ruleID, rule := range rules {
		if rule.Match == "" {
			logger.Warn().
				Str("log_id", logID).
				Int("rule_id", ruleID).
				Msg("invalid metric rule, empty 'match', skipping config")
			return false
		}

		if rule.Name == "" {
			logger.Warn().
				Str("log_id", logID).
				Int("rule_id", ruleID).
				Msg("invalid metric rule, empty 'name', skipping config")
			return false
		}

		matcher, err := regexp.Compile(rule.Match)
		if err != nil {
			logger.Warn().
				Err(err).
				Str("log_id", logID).
				Int("rule_id", ruleID).
				Str("match", rule.Match).
				Msg("rule match compile failed, skipping config")
			return false
		}
		if matcher == nil {
			logger.Warn().
				Str("log_id", logID).
				Int("rule_id", ruleID).
				Str("match", rule.Match).
				Msg("rule match compile resulted in nil value, skipping config")
			return false
		}
		rule.Matcher = matcher
		rule.MatchParts = matcher.SubexpNames()

		if len(rule.MatchParts) < 2 {
			logger.Warn().
				Str("log_id", logID).
				Int("rule_id", ruleID).
				Msg("forcing type to counter, no named subexpressions found")
			rule.Type = "c"
		} else {
			// find the 'Value' subexpression and save its index for extraction on matched lines
			for _, subName := range rule.MatchParts {
				if strings.ToLower(subName) == "value" {
					rule.ValueKey = subName
					break // there can be only one
				}
			}

			if rule.ValueKey == "" {
				logger.Warn().
					Str("log_id", logID).
					Int("id_id", ruleID).
					Msg("forcing type to counter, no subexpression named 'Value' found")
				rule.Type = "c"
			}
		}

		// name contains template interpolation code
		if strings.Contains(rule.Name, "{{.") {
			if len(rule.MatchParts) < 2 {
				logger.Warn().
					Str("log_id", logID).
					Int("id_id", ruleID).
					Str("match", rule.Match).
					Str("name", rule.Name).
					Msg("'name' expects matches, match has no named subexpressions, skipping config")
				return false
			}
			templateID := fmt.Sprintf("%s:M%d-name", logID, ruleID)
			namer, err := template.New(templateID).Parse(rule.Name)
			if err != nil {
				logger.Warn().
					Err(err).
					Str("log_id", logID).
					Int("id_id", ruleID).
					Str("name", rule.Name).
					Msg("name template parse failed, skipping config")
				return false
			}
			if namer == nil {
				logger.Warn().
					Str("log_id", logID).
					Int("id_id", ruleID).
					Str("name", rule.Name).
					Msg("name template parse resulted in nil value, skipping config")
				return false
			}
			rule.Namer = namer
		}

		// tags contains template interpolation code
		if strings.Contains(rule.Tags, "{{.") {
			if len(rule.MatchParts) < 2 {
				logger.Warn().
					Str("log_id", logID).
					Int("id_id", ruleID).
					Str("match", rule.Match).
					Str("tags", rule.Tags).
					Msg("'tags' expects matches, match has no named subexpressions, skipping config")
				return false
			}
			templateID := fmt.Sprintf("%s:M%d-tags", logID, ruleID)
			tagger, err := template.New(templateID).Parse(rule.Tags)
			if err != nil {
				logger.Warn().
					Err(err).
					Str("log_id", logID).
					Int("id_id", ruleID).
					Str("tags", rule.Tags).
					Msg("tags template parse failed, skipping config")
				return false
			}
			if tagger == nil {
				logger.Warn().
					Str("log_id", logID).
					Int("id_id", ruleID).
					Str("tags", rule.Tags).
					Msg("tags template parse resulted in nil value, skipping config")
				return false
			}
			rule.Tagger = tagger
		}

	}

	return true
}

// parse reads and parses a log configuration
func parse(cfgType, cfgFile string) (Config, error) {
	var cfg Config

	data, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		return cfg, err
	}

	switch cfgType {
	case ".json":
		err := json.Unmarshal(data, &cfg)
		if serr, ok := err.(*json.SyntaxError); ok {
			line, col := findLine(data, serr.Offset)
			return cfg, errors.Wrapf(err, "line %d, col %d", line, col)
		}
		return cfg, err
	case ".yaml":
		err := yaml.Unmarshal(data, &cfg)
		if err != nil {
			return cfg, err
		}
	case ".toml":
		err := toml.Unmarshal(data, &cfg)
		if err != nil {
			return cfg, err
		}
	default:
		return cfg, errors.Errorf("unknown config type (%s)", cfgType)
	}

	return cfg, nil
}

// checkLogFileAccess verifies a log file can be opened for reading
func checkLogFileAccess(file string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	f.Close()
	return nil
}

func findLine(data []byte, offset int64) (int, int) {
	if offset == 0 {
		return 1, 1
	}
	const (
		CR = 13
		LF = 10
	)
	line := 1
	col := 0
	for i := int64(0); i < offset; i++ {
		col++
		if data[i] == CR || data[i] == LF {
			col = 0
			line++
		}
	}
	return line, col
}
