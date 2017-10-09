// Copyright © 2017 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package logonly

import (
	"sync"

	"github.com/rs/zerolog"
)

// LogOnly defines logging metrics only destination
type LogOnly struct {
	logger zerolog.Logger
}

var (
	client *LogOnly
	once   sync.Once
)
