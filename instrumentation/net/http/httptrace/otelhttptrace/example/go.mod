module go.opentelemetry.io/contrib/instrumentation/net/http/httptrace/otelhttptrace/example

go 1.14

replace (
	go.opentelemetry.io/contrib => ../../../../../../
	go.opentelemetry.io/contrib/instrumentation/net/http/httptrace/otelhttptrace => ../
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp => ../../../otelhttp
)

require (
	go.opentelemetry.io/contrib/instrumentation/net/http/httptrace/otelhttptrace v0.12.0
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.12.0
	go.opentelemetry.io/otel/exporters/stdout v0.12.0
	go.opentelemetry.io/otel/sdk v1.4.1
	go.opentelemetry.io/otel/sdk/export/metric v0.27.0 // indirect
)
