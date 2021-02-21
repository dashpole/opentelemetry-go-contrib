module go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux

go 1.14

replace (
	go.opentelemetry.io/contrib => ../../../../..
	go.opentelemetry.io/contrib/propagators => ../../../../../propagators
)

require (
	github.com/DataDog/sketches-go v0.0.1 // indirect
	github.com/google/gofuzz v1.1.0 // indirect
	github.com/gorilla/mux v1.8.0
	github.com/stretchr/testify v1.7.0
	go.opentelemetry.io/contrib v0.12.0
	go.opentelemetry.io/contrib/propagators v0.12.0
	go.opentelemetry.io/otel v0.17.0
	go.opentelemetry.io/otel/exporters/stdout v0.17.0
	go.opentelemetry.io/otel/sdk v0.17.0
	google.golang.org/grpc v1.32.0 // indirect
)
