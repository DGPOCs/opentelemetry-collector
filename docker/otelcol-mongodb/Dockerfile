# syntax=docker/dockerfile:1
FROM gcr.io/distroless/base-debian12

COPY otelcorecol /otelcol
COPY otelcol.yaml /etc/otelcol/config.yaml

USER 0
ENTRYPOINT ["/otelcol", "--config=/etc/otelcol/config.yaml"]
