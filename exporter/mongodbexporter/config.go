// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package mongodbexporter // import "go.opentelemetry.io/collector/exporter/mongodbexporter"

import (
	"fmt"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configopaque"
	"go.opentelemetry.io/collector/config/configretry"
	"go.opentelemetry.io/collector/config/configtls"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
)

// Config defines the configuration for the MongoDB exporter.
type Config struct {
	exporterhelper.TimeoutConfig `mapstructure:",squash"`
	QueueConfig                  exporterhelper.QueueBatchConfig `mapstructure:"sending_queue"`
	RetryConfig                  configretry.BackOffConfig       `mapstructure:"retry_on_failure"`

	// URI is the MongoDB connection URI. It must include the protocol (e.g. mongodb or mongodb+srv).
	URI string `mapstructure:"uri"`

	// Username is the optional username used when the URI does not embed credentials.
	Username string `mapstructure:"username"`

	// Password is the optional password used together with Username.
	Password configopaque.String `mapstructure:"password"`

	// Database is the MongoDB database that will receive telemetry documents.
	Database string `mapstructure:"database"`

	// Collection is the MongoDB collection that will receive telemetry documents. All signal types share
	// the same collection. Users that prefer different collections can configure multiple exporters with
	// different collection names.
	Collection string `mapstructure:"collection"`

	// TLS controls optional TLS configuration when connecting to MongoDB.
	TLS configtls.ClientConfig `mapstructure:"tls"`

	// ConnectTimeout defines how long the exporter waits for the initial connection and for shutdown
	// disconnect operations. Defaults to 10s when unset.
	ConnectTimeout time.Duration `mapstructure:"connect_timeout"`

	// prevent unkeyed literal initialization
	_ struct{}
}

var _ component.Config = (*Config)(nil)

const defaultConnectionTimeout = 10 * time.Second

// Validate checks if the configuration is valid.
func (c *Config) Validate() error {
	if c.URI == "" {
		return fmt.Errorf("mongodb exporter requires a non-empty \"uri\"")
	}
	if c.Database == "" {
		return fmt.Errorf("mongodb exporter requires a non-empty \"database\"")
	}
	if c.Collection == "" {
		return fmt.Errorf("mongodb exporter requires a non-empty \"collection\"")
	}
	return nil
}

// connectionTimeout returns the configured timeout or a sensible default.
func (c *Config) connectionTimeout() time.Duration {
	if c.ConnectTimeout <= 0 {
		return defaultConnectionTimeout
	}
	return c.ConnectTimeout
}
