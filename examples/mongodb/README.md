# MongoDB Exporter Example

This directory contains an example configuration that wires the MongoDB exporter into the default OTLP pipelines. Environment variables referenced from the YAML allow the Collector to obtain credentials at runtime:

- `MONGODB_URI`
- `MONGODB_USER`
- `MONGODB_PASS`
- `MONGODB_DATABASE`
- `MONGODB_COLLECTION`

Combine this configuration with the distribution built from `cmd/otelcorecol` to create a Collector image that persists telemetry into MongoDB. The `docker/otelcol-mongodb/Dockerfile` file shows a minimal container image that copies the custom Collector binary and this configuration file.
