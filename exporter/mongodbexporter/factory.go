// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package mongodbexporter // import "go.opentelemetry.io/collector/exporter/mongodbexporter"

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configretry"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
	"go.opentelemetry.io/collector/exporter/mongodbexporter/internal/metadata"
	"go.opentelemetry.io/collector/exporter/xexporter"
)

// NewFactory creates a factory for the MongoDB exporter.
func NewFactory() exporter.Factory {
	return xexporter.NewFactory(
		metadata.Type,
		createDefaultConfig,
		xexporter.WithTraces(createTracesExporter, metadata.TracesStability),
		xexporter.WithMetrics(createMetricsExporter, metadata.MetricsStability),
		xexporter.WithLogs(createLogsExporter, metadata.LogsStability),
	)
}

func createDefaultConfig() component.Config {
	return &Config{
		TimeoutConfig: exporterhelper.NewDefaultTimeoutConfig(),
		QueueConfig:   exporterhelper.NewDefaultQueueConfig(),
		RetryConfig:   configretry.NewDefaultBackOffConfig(),
	}
}

func createTracesExporter(ctx context.Context, set exporter.Settings, cfg component.Config) (exporter.Traces, error) {
	c := cfg.(*Config)
	exp, err := newMongoExporter(c, set)
	if err != nil {
		return nil, err
	}
	return exporterhelper.NewTraces(ctx, set, cfg,
		exp.pushTraces,
		exporterhelper.WithTimeout(c.TimeoutConfig),
		exporterhelper.WithQueue(c.QueueConfig),
		exporterhelper.WithRetry(c.RetryConfig),
		exporterhelper.WithStart(exp.start),
		exporterhelper.WithShutdown(exp.shutdown),
	)
}

func createMetricsExporter(ctx context.Context, set exporter.Settings, cfg component.Config) (exporter.Metrics, error) {
	c := cfg.(*Config)
	exp, err := newMongoExporter(c, set)
	if err != nil {
		return nil, err
	}
	return exporterhelper.NewMetrics(ctx, set, cfg,
		exp.pushMetrics,
		exporterhelper.WithTimeout(c.TimeoutConfig),
		exporterhelper.WithQueue(c.QueueConfig),
		exporterhelper.WithRetry(c.RetryConfig),
		exporterhelper.WithStart(exp.start),
		exporterhelper.WithShutdown(exp.shutdown),
	)
}

func createLogsExporter(ctx context.Context, set exporter.Settings, cfg component.Config) (exporter.Logs, error) {
	c := cfg.(*Config)
	exp, err := newMongoExporter(c, set)
	if err != nil {
		return nil, err
	}
	return exporterhelper.NewLogs(ctx, set, cfg,
		exp.pushLogs,
		exporterhelper.WithTimeout(c.TimeoutConfig),
		exporterhelper.WithQueue(c.QueueConfig),
		exporterhelper.WithRetry(c.RetryConfig),
		exporterhelper.WithStart(exp.start),
		exporterhelper.WithShutdown(exp.shutdown),
	)
}
