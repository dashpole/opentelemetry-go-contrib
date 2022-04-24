module go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc/example

go 1.14

replace go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc => ../

require (
	github.com/golang/protobuf v1.5.2
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.12.0
	go.opentelemetry.io/otel v0.12.0
	go.opentelemetry.io/otel/exporters/stdout v0.12.0
	go.opentelemetry.io/otel/sdk v0.12.0
	golang.org/x/net v0.0.0-20201021035429-f5854403a974
	google.golang.org/grpc v1.46.0
)
