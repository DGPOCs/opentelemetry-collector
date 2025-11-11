// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package metadata

import "go.opentelemetry.io/collector/component"

var (
	Type      = component.MustNewType("mongodb")
	ScopeName = "go.opentelemetry.io/collector/exporter/mongodbexporter"
)

const (
	TracesStability  = component.StabilityLevelAlpha
	MetricsStability = component.StabilityLevelAlpha
	LogsStability    = component.StabilityLevelAlpha
)
