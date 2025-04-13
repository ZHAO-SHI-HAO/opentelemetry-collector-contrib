package duckdbexporter

import (
	"database/sql"
	"errors"
	_ "github.com/marcboeker/go-duckdb/v2"
	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/duckdbexporter/internal"
	"go.opentelemetry.io/collector/config/configretry"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
)

type Config struct {
	TimeoutSettings           exporterhelper.TimeoutConfig `mapstructure:",squash"`
	configretry.BackOffConfig `mapstructure:"retry_on_failure"`
	QueueSettings             exporterhelper.QueueBatchConfig `mapstructure:"sending_queue"`

	// Database is the database name to export.
	Database string `mapstructure:"database"`
	// ConnectionParams is the extra connection parameters with map format. for example compression/dial_timeout
	ConnectionParams map[string]string `mapstructure:"connection_params"`
	// LogsTableName is the table name for logs. default is `otel_logs`.
	LogsTableName string `mapstructure:"logs_table_name"`
	// TracesTableName is the table name for traces. default is `otel_traces`.
	TracesTableName string `mapstructure:"traces_table_name"`
	// CreateSchema if set to true will run the DDL for creating the database and tables. default is true.
	CreateSchema bool `mapstructure:"create_schema"`
	// MetricsTables defines the table names for metric types.
	MetricsTables MetricTablesConfig `mapstructure:"metrics_tables"`
}

type MetricTablesConfig struct {
	// Gauge is the table name for gauge metric type. default is `otel_metrics_gauge`.
	Gauge internal.MetricTypeConfig `mapstructure:"gauge"`
	// Sum is the table name for sum metric type. default is `otel_metrics_sum`.
	Sum internal.MetricTypeConfig `mapstructure:"sum"`
	// Summary is the table name for summary metric type. default is `otel_metrics_summary`.
	Summary internal.MetricTypeConfig `mapstructure:"summary"`
	// Histogram is the table name for histogram metric type. default is `otel_metrics_histogram`.
	Histogram internal.MetricTypeConfig `mapstructure:"histogram"`
	// ExponentialHistogram is the table name for exponential histogram metric type. default is `otel_metrics_exponential_histogram`.
	ExponentialHistogram internal.MetricTypeConfig `mapstructure:"exponential_histogram"`
}

const (
	defaultDatabase           = "default"
	defaultMetricTableName    = "otel_metrics"
	defaultGaugeSuffix        = "_gauge"
	defaultSumSuffix          = "_sum"
	defaultSummarySuffix      = "_summary"
	defaultHistogramSuffix    = "_histogram"
	defaultExpHistogramSuffix = "_exponential_histogram"
)

var (
	errConfigNoDataBase = errors.New("database must be specified")
)

// Validate the duckDb configuration.
func (cfg *Config) Validate() (err error) {
	if cfg.Database == "" {
		err = errors.Join(err, errConfigNoDataBase)
	}

	return err
}

func (cfg *Config) buildDB() (*sql.DB, error) {
	conn, err := sql.Open("duckdb", cfg.Database)
	defer conn.Close()
	if err != nil {
		return nil, err
	}

	return conn, nil
}

// shouldCreateSchema returns true if the exporter should run the DDL for creating database/tables.
func (cfg *Config) shouldCreateSchema() bool {
	return cfg.CreateSchema
}

func (cfg *Config) buildMetricTableNames() {
	if len(cfg.MetricsTables.Gauge.Name) == 0 {
		cfg.MetricsTables.Gauge.Name = defaultGaugeSuffix
	}
	if len(cfg.MetricsTables.Sum.Name) == 0 {
		cfg.MetricsTables.Sum.Name = defaultSumSuffix
	}
	if len(cfg.MetricsTables.Summary.Name) == 0 {
		cfg.MetricsTables.Summary.Name = defaultSummarySuffix
	}
	if len(cfg.MetricsTables.Histogram.Name) == 0 {
		cfg.MetricsTables.Histogram.Name = defaultHistogramSuffix
	}
	if len(cfg.MetricsTables.ExponentialHistogram.Name) == 0 {
		cfg.MetricsTables.ExponentialHistogram.Name = defaultExpHistogramSuffix
	}
}
