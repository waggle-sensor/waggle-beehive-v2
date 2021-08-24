FROM golang:1.17 AS builder
WORKDIR /build
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o influxdb-data-api

FROM scratch
COPY --from=builder /build/influxdb-data-api /influxdb-data-api
ENTRYPOINT [ "/influxdb-data-api" ]
