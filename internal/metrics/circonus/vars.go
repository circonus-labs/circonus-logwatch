// Copyright © 2017 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package circonus

import (
	cgm "github.com/circonus-labs/circonus-gometrics"
	"github.com/rs/zerolog"
	tomb "gopkg.in/tomb.v2"
)

// Circonus defines an instance of the circonus metrics destination
type Circonus struct {
	logger zerolog.Logger
	client *cgm.CirconusMetrics
	t      tomb.Tomb
}
