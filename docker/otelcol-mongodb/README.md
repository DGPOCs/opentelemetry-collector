# Docker Image for the MongoDB Collector

1. Generate the distribution that includes the MongoDB exporter:

   ```bash
   go run go.opentelemetry.io/collector/cmd/builder --config=cmd/otelcorecol/builder-config.yaml --output-path=./dist
   ```

2. Copy the generated binary and configuration into the Docker build context:

   ```bash
   cp dist/otelcorecol ./docker/otelcol-mongodb/
   cp examples/mongodb/otelcol.yaml ./docker/otelcol-mongodb/
   ```

3. Build the image:

   ```bash
   docker build -t myorg/otelcol-mongodb:latest docker/otelcol-mongodb
   ```

4. Reference the image from `docker-compose.yaml` or another orchestrator, providing the MongoDB settings via environment variables.
