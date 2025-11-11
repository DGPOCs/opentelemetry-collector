// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package mongodbexporter // import "go.opentelemetry.io/collector/exporter/mongodbexporter"

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

type signalType string

const (
	signalTraces  signalType = "traces"
	signalMetrics signalType = "metrics"
	signalLogs    signalType = "logs"
)

type mongoExporter struct {
	cfg        *Config
	client     *mongo.Client
	collection *mongo.Collection
	logger     *zap.Logger
	tracesEnc  ptrace.Marshaler
	metricsEnc pmetric.Marshaler
	logsEnc    plog.Marshaler
}

func newMongoExporter(cfg *Config, set exporter.Settings) (*mongoExporter, error) {
	clientOpts := options.Client().ApplyURI(cfg.URI)
	if cfg.Username != "" {
		clientOpts.SetAuth(options.Credential{Username: cfg.Username, Password: string(cfg.Password)})
	}

	if cfg.TLS.Insecure && cfg.TLS.InsecureSkipVerify {
		return nil, fmt.Errorf("mongodb exporter cannot set both tls.insecure and tls.insecure_skip_verify")
	}

	tlsCfg, err := cfg.TLS.LoadTLSConfig(context.Background())
	if err != nil {
		return nil, err
	}
	if tlsCfg != nil {
		clientOpts.SetTLSConfig(tlsCfg)
	}

	client, err := mongo.NewClient(clientOpts)
	if err != nil {
		return nil, err
	}

	logger := set.TelemetrySettings.Logger
	if logger == nil {
		logger = zap.NewNop()
	}

	return &mongoExporter{
		cfg:        cfg,
		client:     client,
		logger:     logger,
		tracesEnc:  &ptrace.ProtoMarshaler{},
		metricsEnc: &pmetric.ProtoMarshaler{},
		logsEnc:    &plog.ProtoMarshaler{},
	}, nil
}

func (e *mongoExporter) start(ctx context.Context, _ component.Host) error {
	connectCtx, cancel := context.WithTimeout(ctx, e.cfg.connectionTimeout())
	defer cancel()
	if err := e.client.Connect(connectCtx); err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %w", err)
	}
	e.collection = e.client.Database(e.cfg.Database).Collection(e.cfg.Collection)
	return nil
}

func (e *mongoExporter) shutdown(ctx context.Context) error {
	disconnectCtx, cancel := context.WithTimeout(ctx, e.cfg.connectionTimeout())
	defer cancel()
	return e.client.Disconnect(disconnectCtx)
}

func (e *mongoExporter) pushTraces(ctx context.Context, td ptrace.Traces) error {
	if td.SpanCount() == 0 {
		return nil
	}
	bytes, err := e.tracesEnc.MarshalTraces(td)
	if err != nil {
		return err
	}
	return e.insertDocument(ctx, signalTraces, bytes, int64(td.SpanCount()))
}

func (e *mongoExporter) pushMetrics(ctx context.Context, md pmetric.Metrics) error {
	if md.DataPointCount() == 0 {
		return nil
	}
	bytes, err := e.metricsEnc.MarshalMetrics(md)
	if err != nil {
		return err
	}
	return e.insertDocument(ctx, signalMetrics, bytes, int64(md.DataPointCount()))
}

func (e *mongoExporter) pushLogs(ctx context.Context, ld plog.Logs) error {
	if ld.LogRecordCount() == 0 {
		return nil
	}
	bytes, err := e.logsEnc.MarshalLogs(ld)
	if err != nil {
		return err
	}
	return e.insertDocument(ctx, signalLogs, bytes, int64(ld.LogRecordCount()))
}

func (e *mongoExporter) insertDocument(ctx context.Context, signal signalType, payload []byte, count int64) error {
	if e.collection == nil {
		return fmt.Errorf("mongodb exporter is not connected")
	}
	doc := bson.M{
		"signal":    signal,
		"encoding":  "proto",
		"payload":   payload,
		"items":     count,
		"timestamp": time.Now().UTC(),
	}
	if _, err := e.collection.InsertOne(ctx, doc); err != nil {
		return err
	}
	e.logger.Debug("MongoDB exporter stored telemetry", zap.String("signal", string(signal)), zap.Int64("items", count))
	return nil
}
