module go.opentelemetry.io/contrib/instrumentation/github.com/bradfitz/gomemcache/memcache/otelmemcache

go 1.14

replace go.opentelemetry.io/contrib => ../../../../../../

require (
	github.com/bradfitz/gomemcache v0.0.0-20190913173617-a41fca850d0b
	github.com/stretchr/testify v1.6.1
	go.opentelemetry.io/contrib v0.22.0
	go.opentelemetry.io/otel v0.12.0
	go.opentelemetry.io/otel/exporters/stdout v0.12.0
	go.opentelemetry.io/otel/sdk v0.12.0
	google.golang.org/genproto v0.0.0-20200331122359-1ee6d9798940 // indirect
)
