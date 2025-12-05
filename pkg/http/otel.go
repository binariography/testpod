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

	"github.com/binariography/testpod/pkg/version"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/contrib/propagators/aws/xray"
	"go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/contrib/propagators/jaeger"
	"go.opentelemetry.io/contrib/propagators/ot"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
)

var (
	instrumentationName = "github.com/binariography/testpod/pkg/http"
)

func (s *Server) initTracer(ctx context.Context) {
	if s.config.OtelService == "" {
		noop := noop.NewTracerProvider()
		s.tracer = noop.Tracer("")
		return
	}

	traceExporter, _ := otlptracehttp.New(ctx)
	s.tracerProvider = sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(s.config.OtelService),
			semconv.ServiceVersionKey.String(version.VERSION),
		)))

	otel.SetTracerProvider(s.tracerProvider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
		b3.New(),
		&jaeger.Jaeger{},
		&ot.OT{},
		&xray.Propagator{},
	))

	s.tracer = s.tracerProvider.Tracer(
		instrumentationName,
		trace.WithInstrumentationVersion(version.VERSION),
		trace.WithSchemaURL(semconv.SchemaURL),
	)

}

func NewOtelMiddleware(srv string) mux.MiddlewareFunc {
	return otelhttp.NewMiddleware(srv)
}
