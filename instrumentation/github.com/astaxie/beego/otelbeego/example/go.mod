module go.opentelemetry.io/contrib/instrumentation/astaxie/beego/example

go 1.14

require (
	github.com/astaxie/beego v1.12.2
	go.opentelemetry.io/contrib/instrumentation/github.com/astaxie/beego/otelbeego v0.12.0
	go.opentelemetry.io/otel/exporters/stdout v0.12.0
	go.opentelemetry.io/otel/sdk v1.9.0
	go.opentelemetry.io/otel/sdk/export/metric v0.28.0 // indirect
	go.opentelemetry.io/otel/sdk/metric v0.31.0 // indirect
)

replace (
	go.opentelemetry.io/contrib => ../../../../../../
	go.opentelemetry.io/contrib/instrumentation/github.com/astaxie/beego/otelbeego => ../
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp => ../../../../../net/http/otelhttp
	go.opentelemetry.io/contrib/propagators => ../../../../../../propagators
)
