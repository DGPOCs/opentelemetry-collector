# MongoDB Exporter

The MongoDB exporter writes telemetry payloads from the OpenTelemetry Collector into a MongoDB collection. Each export operation stores the serialized OTLP payload as a binary field together with metadata describing the signal type and the number of items contained in the batch.

## Configuration

```yaml
exporters:
  mongodb:
    uri: "${env:MONGODB_URI}"
    username: "${env:MONGODB_USER}"
    password: "${env:MONGODB_PASS}"
    database: telemetry
    collection: signals
    connect_timeout: 15s
    tls:
      insecure: false
```

Use the exporter in a pipeline by referencing it from the desired signal type:

```yaml
service:
  pipelines:
    traces:
      receivers: [otlp]
      processors: [batch]
      exporters: [mongodb]
```

All sensitive values can be provided via environment variables using the Collector's standard `${env:...}` notation.
