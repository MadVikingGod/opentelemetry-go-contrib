module go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-sdk-go/aws/otelaws/example

go 1.15

replace (
	// go.opentelemetry.io/contrib => ../../../../../../../
	go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-sdk-go/aws/otelaws => ../
	go.opentelemetry.io/contrib/propagators => ../../../../../../../propagators
)

require (
	github.com/aws/aws-sdk-go v1.36.14
	go.opentelemetry.io/otel v0.15.0
	go.opentelemetry.io/otel/exporters/stdout v0.15.0
	go.opentelemetry.io/otel/sdk v0.15.0
	go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-sdk-go/aws/otelaws v0.15.1
)
