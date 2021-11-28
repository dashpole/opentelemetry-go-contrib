module go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc/example

go 1.14

replace go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc => ../

require (
	github.com/golang/protobuf v1.4.2
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.12.0
	go.opentelemetry.io/otel/exporters/stdout v0.12.0
	go.opentelemetry.io/otel/sdk v1.2.0
	go.opentelemetry.io/otel/sdk/metric v0.25.0 // indirect
	golang.org/x/net v0.0.0-20200707034311-ab3426394381
	google.golang.org/grpc v1.32.0
)
