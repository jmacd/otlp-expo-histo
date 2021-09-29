// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package histogram

import (
	"math"
)

// LogarithmMapping is a prototype for OTEP 149.  The Go
// implementation was copied from a Java prototypes during following
// https://github.com/open-telemetry/opentelemetry-proto/pull/322.
// See
// https://github.com/newrelic-experimental/newrelic-sketch-java/blob/1ce245713603d61ba3a4510f6df930a5479cd3f6/src/main/java/com/newrelic/nrsketch/indexer/LogIndexer.java
// for the equations used here.  Note there are even more options
// implemented in that package.
type LogarithmMapping struct {
	scale int

	// scaleFactor is used and computed as follows:
	// index = log(value) / log(base)
	// = log(value) / log(2^(2^-scale))
	// = log(value) / (2^-scale * log(2))
	// = log(value) * (1/log(2) * 2^scale)
	// = log(value) * scaleFactor
	// where:
	// scaleFactor = (1/log(2) * 2^scale)
	// = math.Log2E * math.Exp2(scale)
	// = scalb(math.Log2E, scale)
	// Because multiplication is faster than division, we define scaleFactor as a multiplier.
	scaleFactor float64
}

var _ Base2HistogramMapper = &LogarithmMapping{}

func NewLogarithmMapping(scale int) *LogarithmMapping {
	return &LogarithmMapping{
		scale:       scale,
		scaleFactor: Scalb(math.Log2E, scale),
	}
}

func (l *LogarithmMapping) MapToIndex(value float64) int64 {
	// Use Floor() to round toward -Inf.
	return int64(math.Floor(math.Log(value) * l.scaleFactor))
}

func (l *LogarithmMapping) LowerBoundary(index int64) float64 {
	// result = base ^ index
	// = (2^(2^-scale))^index
	// = 2^(2^-scale * index)
	// = 2^(index * 2^-scale))
	return math.Exp2(Scalb(float64(index), -l.scale))
}

func (l *LogarithmMapping) UpperBoundary(index int64) float64 {
	return l.LowerBoundary(index + 1)
}
