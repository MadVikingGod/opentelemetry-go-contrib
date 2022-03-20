module go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp/example

go 1.14

replace (
	go.opentelemetry.io/contrib => ../../../../../
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp => ../
)

require (
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.15.1
	go.opentelemetry.io/otel v1.5.0
	go.opentelemetry.io/otel/exporters/stdout v0.15.0
	go.opentelemetry.io/otel/oteltest v0.20.0 // indirect
	go.opentelemetry.io/otel/sdk v1.5.0
	go.opentelemetry.io/otel/sdk/export/metric v0.27.0 // indirect
	go.opentelemetry.io/otel/trace v1.5.0
)
