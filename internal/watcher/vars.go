// Copyright © 2017 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package watcher

import (
	"github.com/circonus-labs/circonus-logwatch/internal/configs"
	"github.com/circonus-labs/circonus-logwatch/internal/metrics"
	"github.com/rs/zerolog"
	tomb "gopkg.in/tomb.v2"
)

type metric struct {
	Name  string
	Type  string
	Value string
}

type metricLine struct {
	line     string
	matches  *map[string]string
	metricID int
}

// Watcher defines a new log watcher
type Watcher struct {
	cfg              *configs.Config
	trace            bool
	logger           zerolog.Logger
	metricLines      chan metricLine
	metrics          chan metric
	t                tomb.Tomb
	dest             metrics.Destination
	statTotalLines   string
	statMatchedLines string
}

const (
	metricLineQueueSize = 1000
	metricQueueSize     = 1000
)
