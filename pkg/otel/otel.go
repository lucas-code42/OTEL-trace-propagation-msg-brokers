package opentelemetry

import (
	"context"
	"log"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"
)

type opentelemetrySpan struct {
	span trace.Span
}

type openTelemetryTrace struct {
	trace trace.Tracer
}

func New(ctx context.Context, serviceName string) *openTelemetryTrace {
	return &openTelemetryTrace{
		trace: otelInstrumentarion(ctx, serviceName),
	}
}

func otelInstrumentarion(ctx context.Context, serviceName string) trace.Tracer {
	client := otlptracegrpc.NewClient(otlptracegrpc.WithInsecure())

	exporter, err := otlptrace.New(ctx, client)
	if err != nil {
		log.Fatalf("%s:", err)
	}

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(newResource(serviceName)),
	)
	return tracerProvider.Tracer(serviceName, trace.WithInstrumentationVersion("1.0"))
}

func newResource(service string) *resource.Resource {
	return resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName(service),
		semconv.ServiceVersion("0.0.1"),
	)
}

func (o *openTelemetryTrace) Start(ctx context.Context, spanName string) (context.Context, *opentelemetrySpan) {
	ctx, span := o.trace.Start(ctx, spanName)
	return ctx, &opentelemetrySpan{span: span}
}

func (s *opentelemetrySpan) End() {
	s.span.End()
}
