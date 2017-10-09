// Copyright © 2017 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package statsd

import (
	"net"
	"sync"

	"github.com/rs/zerolog"
)

// Statsd defines the relevant properties of a StatsD connection.
type Statsd struct {
	id     string
	port   string
	prefix string
	conn   net.Conn
	logger zerolog.Logger
}

var (
	client *Statsd
	once   sync.Once
)
