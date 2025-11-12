/*
Docs

Traces give us the big picture of what happens when a request is made to an application. Whether your application is a monolith with a single database or a sophisticated mesh of services, traces are essential to understanding the full “path” a request takes in your application.

A span represents a unit of work or operation. Spans are the building blocks of Traces. In OpenTelemetry, they include the following information:

    Name
    Parent span ID (empty for root spans)
    Start and End Timestamps
    Span Context
    Attributes
    Span Events
    Span Links
    Span Status


handleFunc := func(pattern string,
	handlerFunc func(http.ResponseWriter, *http.Request)) {
                // Configure the "http.route" for the HTTP instrumentation.
                handler := otelhttp.WithRouteTag(pattern, http.HandlerFunc(handlerFunc))
                mux.Handle(pattern, handler)
        }


func NewMiddleware(operation string, opts ...Option) func(http.Handler) http.Handler








*/

package http

import (
	"context"

	"github.com/gorilla/mux"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"

	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
)

var (
	instrumentationName = "testpod"
)

func (s *Server) initTracer(ctx context.Context) {
	//traceExporter, _ := stdouttrace.New(
	//	stdouttrace.WithPrettyPrint())
	traceExporter, _ := stdouttrace.New()

	s.tracerProvider = sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExporter),
	)

	otel.SetTracerProvider(s.tracerProvider)

	s.tracer = s.tracerProvider.Tracer(
		instrumentationName,
		trace.WithInstrumentationVersion("1.0.1"),
		trace.WithSchemaURL(semconv.SchemaURL),
	)

}

func NewOtelMiddleware() mux.MiddlewareFunc {
	return otelmux.Middleware("testpod-front")
}
