module go.opentelemetry.io/contrib/instrumentation/host

go 1.14

replace go.opentelemetry.io/contrib => ../..

require (
	github.com/shirou/gopsutil v2.21.11+incompatible
	github.com/stretchr/testify v1.6.1
	github.com/tklauser/go-sysconf v0.3.9 // indirect
	github.com/yusufpapurcu/wmi v1.2.2 // indirect
	go.opentelemetry.io/contrib v0.12.0
	go.opentelemetry.io/otel v0.12.0
	go.opentelemetry.io/otel/exporters/stdout v0.12.0
	go.opentelemetry.io/otel/sdk v0.12.0
)
