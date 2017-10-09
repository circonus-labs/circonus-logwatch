// Copyright © 2017 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package agent

import (
	"net/http"
	"os"

	"github.com/circonus-labs/circonus-logwatch/internal/metrics"
	"github.com/circonus-labs/circonus-logwatch/internal/watcher"
	tomb "gopkg.in/tomb.v2"
)

// Agent holds the main circonus-logwatch process
type Agent struct {
	watchers   []*watcher.Watcher
	signalCh   chan os.Signal
	t          tomb.Tomb
	destClient metrics.Destination
	svrHTTP    *http.Server
}
