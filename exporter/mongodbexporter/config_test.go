// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package mongodbexporter

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestConfigValidate(t *testing.T) {
	cfg := createDefaultConfig().(*Config)
	require.Error(t, cfg.Validate())

	cfg.URI = "mongodb://localhost:27017"
	require.Error(t, cfg.Validate())

	cfg.Database = "otel"
	require.Error(t, cfg.Validate())

	cfg.Collection = "signals"
	require.NoError(t, cfg.Validate())
}

func TestConnectionTimeoutDefault(t *testing.T) {
	cfg := createDefaultConfig().(*Config)
	require.Equal(t, defaultConnectionTimeout, cfg.connectionTimeout())

	cfg.ConnectTimeout = -1
	require.Equal(t, defaultConnectionTimeout, cfg.connectionTimeout())

	cfg.ConnectTimeout = 3 * time.Second
	require.Equal(t, 3*time.Second, cfg.connectionTimeout())
}
