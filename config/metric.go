// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package config // import "go.opentelemetry.io/contrib/config"

import (
	"context"
	"errors"

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/noop"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
)

func meterProvider(cfg configOptions, res *resource.Resource) (metric.MeterProvider, shutdownFunc, error) {
	if cfg.opentelemetryConfig.MeterProvider == nil {
		return noop.NewMeterProvider(), noopShutdown, nil
	}
	opts := []sdkmetric.Option{
		sdkmetric.WithResource(res),
	}

	var errs []error
	for _, reader := range cfg.opentelemetryConfig.MeterProvider.Readers {
		r, err := metricReader(cfg.ctx, reader)
		if err == nil {
			opts = append(opts, sdkmetric.WithReader(r))
		} else {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return noop.NewMeterProvider(), noopShutdown, errors.Join(errs...)
	}

	mp := sdkmetric.NewMeterProvider(opts...)
	return mp, mp.Shutdown, nil
}

func metricReader(ctx context.Context, r MetricReader) (sdkmetric.Reader, error) {
	if r.Periodic != nil && r.Pull != nil {
		return nil, errors.New("must not specify multiple metric reader type")
	}

	if r.Periodic != nil {
		exp, err := metricExporter(ctx, r.Periodic.Exporter)
		if err != nil {
			return nil, err
		}
		return periodMetricReader(r.Periodic, exp)
	}

	if r.Pull != nil {

	}
	return nil, errors.New("no valid metric reader")
}

func periodMetricReader(pmr *PeriodicMetricReader, exp sdkmetric.Exporter) (sdkmetric.Reader, error) {
	var opts []sdkmetric.PeriodicReaderOption
	return sdkmetric.NewPeriodicReader(exp, opts...), nil
}

func metricExporter(ctx context.Context, exporter MetricExporter) (sdkmetric.Exporter, error) {
	if exporter.Console != nil && exporter.OTLP != nil {
		return nil, errors.New("must not specify multiple exporters")
	}
	return nil, errors.New("no valid metric exporter")
}
