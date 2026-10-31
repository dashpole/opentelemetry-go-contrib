module go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp

go 1.26.0

require (
	github.com/felixge/httpsnoop v1.1.0
	github.com/stretchr/testify v1.12.1
	go.opentelemetry.io/otel v1.47.0
	go.opentelemetry.io/otel/metric v1.47.0
	go.opentelemetry.io/otel/metric/x v0.69.0
	go.opentelemetry.io/otel/sdk v1.47.0
	go.opentelemetry.io/otel/sdk/metric v1.47.0
	go.opentelemetry.io/otel/sdk/metric/x v0.69.0
	go.opentelemetry.io/otel/trace v1.47.0
)

replace (
	go.opentelemetry.io/otel => github.com/dashpole/opentelemetry-go v0.11.1-0.20261031122229-05074a4541e5
	go.opentelemetry.io/otel/log => github.com/dashpole/opentelemetry-go/log v0.0.0-20261031122229-05074a4541e5
	go.opentelemetry.io/otel/metric => github.com/dashpole/opentelemetry-go/metric v0.0.0-20261031122229-05074a4541e5
	go.opentelemetry.io/otel/metric/x => github.com/dashpole/opentelemetry-go/metric/x v0.0.0-20261031122229-05074a4541e5
	go.opentelemetry.io/otel/sdk => github.com/dashpole/opentelemetry-go/sdk v0.11.1-0.20261031122229-05074a4541e5
	go.opentelemetry.io/otel/sdk/metric => github.com/dashpole/opentelemetry-go/sdk/metric v0.0.0-20261031122229-05074a4541e5
	go.opentelemetry.io/otel/sdk/metric/x => github.com/dashpole/opentelemetry-go/sdk/metric/x v0.0.0-20261031122229-05074a4541e5
	go.opentelemetry.io/otel/trace => github.com/dashpole/opentelemetry-go/trace v0.0.0-20261031122229-05074a4541e5
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/go-logr/logr v1.4.4 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/google/uuid v1.6.0 // indirect
	go.opentelemetry.io/auto/sdk v1.2.1 // indirect
	go.opentelemetry.io/otel/log v1.47.0 // indirect
	go.yaml.in/yaml/v3 v3.0.5 // indirect
	golang.org/x/sys v0.49.0 // indirect
)
