package tracing

import (
	"context"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
)

// NewOtlpGrpcProvider returns a new gRPC-based OTLP tracing provider.
func NewOtlpGrpcProvider(ctx context.Context, url string, serviceName string) (*trace.TracerProvider, error) {
	exp, err := otlptracegrpc.New(ctx, otlptracegrpc.WithEndpoint(url))
	if err != nil {
		return nil, err
	}

	// Create resource with proper schema version and attributes
	res, err := resource.New(ctx,
		resource.WithSchemaURL(semconv.SchemaURL),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(serviceName),
		),
	)
	if err != nil {
		return nil, err
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exp),
		trace.WithResource(res),
	)

	return tp, nil
}