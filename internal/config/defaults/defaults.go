// Copyright © 2017 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package defaults

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/circonus-labs/circonus-logwatch/internal/release"
)

const (
	// APIURL for circonus
	APIURL = "https://api.circonus.com/v2/"

	// APIApp defines the api app name associated with the api token key
	APIApp = release.NAME

	// Debug is false by default
	Debug = false

	// AppStatPort for accessing runtime metrics (expvar)
	AppStatPort = "33284"

	// LogLevel set to info by default
	LogLevel = "info"

	// LogPretty colored/formatted output to stderr
	LogPretty = false

	// DestinationType where metrics should be sent (agent|check|log|statsd)
	DestinationType = "log"

	// AgentPort for circonus-agent
	AgentPort = "2609"

	// AgentInterval to submit metrics
	AgentInterval = "60s"

	// StatsdPort for circonus-agent
	StatsdPort = "8125"

	// StatsdPrefix to prepend to every metric
	StatsdPrefix = "host."
)

var (
	// BasePath is the "base" directory
	//
	// expected installation structure:
	// base            (e.g. /opt/circonus/logwatch)
	//   /etc          (e.g. /opt/circonus/logwatch/etc)
	//   /etc/log.d    (e.g. /opt/circonus/logwatch/etc/log.d)
	//   /sbin         (e.g. /opt/circonus/logwatch/sbin)
	BasePath = ""

	// EtcPath returns the default etc directory within base directory
	EtcPath = "" // (e.g. /opt/circonus/logwatch/etc)

	// LogConfPath returns the default directory for log configurations within base directory
	LogConfPath = "" // (e.g. /opt/circonus/logwatch/etc/log.d)

	// Target used when destination type is "check"
	Target = ""
)

func init() {
	var exePath string
	var resolvedExePath string
	var err error

	exePath, err = os.Executable()
	if err == nil {
		resolvedExePath, err = filepath.EvalSymlinks(exePath)
		if err == nil {
			BasePath = filepath.Clean(filepath.Join(filepath.Dir(resolvedExePath), "..")) // e.g. /opt/circonus/agent
		}
	}

	if err != nil {
		fmt.Printf("Unable to determine path to binary %v\n", err)
		os.Exit(1)
	}

	EtcPath = filepath.Join(BasePath, "etc")
	LogConfPath = filepath.Join(EtcPath, "log.d")

	Target, err = os.Hostname()
	if err != nil {
		fmt.Printf("Unable to determine hostname for target %v\n", err)
		os.Exit(1)
	}
}
