module github.com/lightstep/otel-launcher-go/lightstep/sdk/metric/example

go 1.18

require (
	github.com/lightstep/otel-launcher-go/lightstep/sdk/metric v1.15.1
	github.com/lightstep/otel-launcher-go/pipelines v1.8.0
	go.opentelemetry.io/proto/otlp v0.19.0
)

require (
	github.com/cenkalti/backoff/v4 v4.2.0 // indirect
	github.com/dgryski/go-farm v0.0.0-20200201041132-a6ae2369ad13 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.7.0 // indirect
	github.com/lightstep/go-expohisto v1.0.0 // indirect
	go.opentelemetry.io/otel v1.14.0 // indirect
	go.opentelemetry.io/otel/metric v0.37.0 // indirect
	go.opentelemetry.io/otel/sdk v1.14.0 // indirect
	go.opentelemetry.io/otel/sdk/metric v0.37.0 // indirect
	go.opentelemetry.io/otel/trace v1.14.0 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.9.0 // indirect
	golang.org/x/net v0.7.0 // indirect
	golang.org/x/sys v0.5.0 // indirect
	golang.org/x/text v0.7.0 // indirect
	google.golang.org/genproto v0.0.0-20230110181048-76db0878b65f // indirect
	google.golang.org/grpc v1.53.0 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
)

replace github.com/lightstep/otel-launcher-go/lightstep/sdk/metric => ../

replace github.com/lightstep/otel-launcher-go/pipelines => ../../../../pipelines
