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

// Destination defines the interface required by the metric destination.
type Destination interface {
	AddSetValue(string, string) error                               // type 's'  - set metric (ala statsd, counts unique values)
	AddSetValueWithTags(string, []string, string) error             // type 's'  - set metric (ala statsd, counts unique values) with tags
	IncrementCounter(string) error                                  // type 'c'  - counter (monotonically increasing value)
	IncrementCounterWithTags(string, []string) error                // type 'c'  - counter (monotonically increasing value) with tags
	IncrementCounterByValue(string, uint64) error                   // type 'c'  - counter (monotonically increasing value)
	IncrementCounterByValueWithTags(string, []string, uint64) error // type 'c'  - counter (monotonically increasing value) with tags
	SetGaugeValue(string, interface{}) error                        // type 'g'  - gauge (ints or floats)
	SetGaugeValueWithTags(string, []string, interface{}) error      // type 'g'  - gauge (ints or floats) with tags
	SetHistogramValue(string, float64) error                        // type 'h'  - histogram
	SetHistogramValueWithTags(string, []string, float64) error      // type 'h'  - histogram with tags
	SetTextValue(string, string) error                              // type 't'  - text metric
	SetTextValueWithTags(string, []string, string) error            // type 't'  - text metric with tags
	SetTimingValue(string, float64) error                           // type 'ms' - histogram
	SetTimingValueWithTags(string, []string, float64) error         // type 'ms' - histogram with tags
	Start() error
	Stop() error
}
