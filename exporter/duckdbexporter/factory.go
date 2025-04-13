//go:generate mdatagen metadata.yaml

package duckdbexporter

import (
	"context"
	"fmt"
	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/duckdbexporter/internal"
	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/duckdbexporter/internal/metadata"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configretry"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
)

// NewFactory creates a factory for ClickHouse exporter.
func NewFactory() exporter.Factory {
	return exporter.NewFactory(
		metadata.Type,
		createDefaultConfig,
		exporter.WithLogs(createLogsExporter, metadata.LogsStability),
		exporter.WithTraces(createTracesExporter, metadata.TracesStability),
		exporter.WithMetrics(createMetricExporter, metadata.MetricsStability),
	)
}

func createDefaultConfig() component.Config {
	return &Config{
		ConnectionParams: map[string]string{},
		Database:         defaultDatabase,
		LogsTableName:    "otel_logs",
		TracesTableName:  "otel_traces",
		CreateSchema:     true,
		MetricsTables: MetricTablesConfig{
			Gauge:                internal.MetricTypeConfig{Name: defaultMetricTableName + defaultGaugeSuffix},
			Sum:                  internal.MetricTypeConfig{Name: defaultMetricTableName + defaultSumSuffix},
			Summary:              internal.MetricTypeConfig{Name: defaultMetricTableName + defaultSummarySuffix},
			Histogram:            internal.MetricTypeConfig{Name: defaultMetricTableName + defaultHistogramSuffix},
			ExponentialHistogram: internal.MetricTypeConfig{Name: defaultMetricTableName + defaultExpHistogramSuffix},
		},
	}
}

// createLogsExporter creates a new exporter for logs.
// Logs are directly inserted into ClickHouse.
func createLogsExporter(
	ctx context.Context,
	set exporter.Settings,
	cfg component.Config,
) (exporter.Logs, error) {
	c := cfg.(*Config)
	exporter, err := newLogsExporter(set.Logger, c)
	if err != nil {
		return nil, fmt.Errorf("cannot configure clickhouse logs exporter: %w", err)
	}

	return exporterhelper.NewLogs(
		ctx,
		set,
		cfg,
		exporter.pushLogsData,
		exporterhelper.WithStart(exporter.start),
		exporterhelper.WithShutdown(exporter.shutdown),
		exporterhelper.WithTimeout(c.TimeoutSettings),
		exporterhelper.WithQueue(c.QueueSettings),
		exporterhelper.WithRetry(c.BackOffConfig),
	)
}
