module go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-lambda-go/otellambda

go 1.16

replace go.opentelemetry.io/contrib => ../../../../..

require (
	github.com/aws/aws-lambda-go v1.29.0
	github.com/stretchr/testify v1.8.0
	go.opentelemetry.io/otel v1.6.3
	go.opentelemetry.io/otel/trace v1.6.3
)
