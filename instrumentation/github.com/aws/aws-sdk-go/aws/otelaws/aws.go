package otelaws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	otelcontrib "go.opentelemetry.io/contrib"
	"go.opentelemetry.io/otel/semconv"
	oteltrace "go.opentelemetry.io/otel/trace"
)

var tracerName = "go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-sdk-go/aws/otelaws"

func WrapSession(s *session.Session, opts ...Option) *session.Session {
	cfg := &config{}
	cfg.apply(opts...)
	s = s.Copy()
	s.Handlers.Send.PushFrontNamed(request.NamedHandler{
		Name: tracerName + ".Send",
		Fn:   sendHandler(cfg),
	})
	s.Handlers.Complete.PushBackNamed(request.NamedHandler{
		Name: tracerName + ".Complete",
		Fn:   completeHandler(cfg),
	})
	return s
}

func sendHandler(cfg *config) func(*request.Request) {
	tracer := cfg.TracerProvider.Tracer(
		tracerName,
		oteltrace.WithInstrumentationVersion(otelcontrib.SemVersion()),
	)

	return func(req *request.Request) {
		ctx := cfg.Propagators.Extract(req.Context(), req.HTTPRequest.Header)
		opts := []oteltrace.SpanOption{
			oteltrace.WithAttributes(semconv.NetAttributesFromHTTPRequest("tcp", req.HTTPRequest)...),
			oteltrace.WithAttributes(semconv.EndUserAttributesFromHTTPRequest(req.HTTPRequest)...),
			oteltrace.WithAttributes(semconv.HTTPServerAttributesFromHTTPRequest(req.ClientInfo.ServiceName, req.Operation.HTTPPath, req.HTTPRequest)...),
			oteltrace.WithSpanKind(oteltrace.SpanKindServer),
		}
		spanName := req.HTTPRequest.URL.Path
		if spanName == "" {
			spanName = fmt.Sprintf("HTTP %s route not found", req.HTTPRequest.Method)
		}
		ctx, _ = tracer.Start(ctx, spanName, opts...)
		req.SetContext(ctx)
	}
}

func completeHandler(cfg *config) func(*request.Request) {
	return func(req *request.Request) {
		span := oteltrace.SpanFromContext(req.Context())
		if req.HTTPResponse != nil {
			status := req.HTTPResponse.StatusCode
			attrs := semconv.HTTPAttributesFromHTTPStatusCode(status)
			spanStatus, spanMessage := semconv.SpanStatusFromHTTPStatusCode(status)

			span.SetAttributes(attrs...)
			span.SetStatus(spanStatus, spanMessage)

		}
		span.End()
	}
}
