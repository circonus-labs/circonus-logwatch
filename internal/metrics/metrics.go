// Copyright © 2017 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

// Package metrics provides an interface for various metric destinations
//
// Metric types supported:
//   c - counter
//   g - gauge
//   h - histogram
//   ms - timing
//   s - set
//   t - text
package metrics

// Destination defines the interface required by the metric destination
type Destination interface {
	AddSetValue(string, string) error             // type 's'  - set metric (ala statsd, counts unique values)
	IncrementCounter(string) error                // type 'c'  - counter (monotonically increasing value)
	IncrementCounterByValue(string, uint64) error // type 'c'  - counter (monotonically increasing value)
	SetGaugeValue(string, interface{}) error      // type 'g'  - gauge (ints or floats)
	SetHistogramValue(string, float64) error      // type 'h'  - histogram
	SetTextValue(string, string) error            // type 't'  - text metric
	SetTimingValue(string, float64) error         // type 'ms' - histogram
	Start() error
	Stop() error
}
