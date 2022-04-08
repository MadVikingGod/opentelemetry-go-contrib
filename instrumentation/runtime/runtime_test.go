// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package runtime_test

import (
	"context"
	goruntime "runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel/sdk/metric/metrictest"
)

func TestRuntime(t *testing.T) {
	err := runtime.Start(
		runtime.WithMinimumReadMemStatsInterval(time.Second),
	)
	assert.NoError(t, err)
	time.Sleep(time.Second)
}

// TODO: Replace with in memory exporter https://github.com/open-telemetry/opentelemetry-go/issues/2722
func getGCCount(exporter *metrictest.Exporter) int64 {
	ex, err := exporter.GetByName("process.runtime.go.gc.count")
	if err != nil {
		panic("Could not locate a process.runtime.go.gc.count metric in test output")
	}
	return ex.Sum.AsInt64()
}

func testMinimumInterval(t *testing.T, shouldHappen bool, opts ...runtime.Option) {
	goruntime.GC()

	var mstats0 goruntime.MemStats
	goruntime.ReadMemStats(&mstats0)
	baseline := int64(mstats0.NumGC)

	provider, exporter := metrictest.NewTestMeterProvider()

	err := runtime.Start(
		append(
			opts,
			runtime.WithMeterProvider(provider),
		)...,
	)
	assert.NoError(t, err)

	goruntime.GC()

	exporter.Collect(context.Background())

	require.EqualValues(t, 1, getGCCount(exporter)-baseline)

	extra := 0
	if shouldHappen {
		extra = 3
	}

	goruntime.GC()
	goruntime.GC()
	goruntime.GC()

	exporter.Collect(context.Background())
	goruntime.ReadMemStats(&mstats0)

	require.EqualValues(t, 1+extra, getGCCount(exporter)-baseline)
}

func TestDefaultMinimumInterval(t *testing.T) {
	testMinimumInterval(t, false)
}

func TestNoMinimumInterval(t *testing.T) {
	testMinimumInterval(t, true, runtime.WithMinimumReadMemStatsInterval(0))
}

func TestExplicitMinimumInterval(t *testing.T) {
	testMinimumInterval(t, false, runtime.WithMinimumReadMemStatsInterval(time.Hour))
}
