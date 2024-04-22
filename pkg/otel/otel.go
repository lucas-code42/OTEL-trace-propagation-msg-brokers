package app

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
	exp, err := newExporter(ctx)
	if err != nil {
		log.Fatalf("failed to initialize exporter: %v", err)
	}

	tp := newTraceProvider(exp, serviceName)
	defer func() { _ = tp.Shutdown(ctx) }()

	return tp.Tracer(serviceName)
}

func newExporter(ctx context.Context) (*otlptrace.Exporter, error) {
	return otlptracegrpc.New(ctx)
}

func newTraceProvider(exp sdktrace.SpanExporter, serviceName string) *sdktrace.TracerProvider {
	r := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName(serviceName),
		semconv.ServiceVersion("0.0.1"),
	)

	return sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(r),
	)
}

func (o *openTelemetryTrace) Start(ctx context.Context, spanName string) (context.Context, *opentelemetrySpan) {
	ctx, span := o.trace.Start(ctx, spanName)
	return ctx, &opentelemetrySpan{span: span}
}
