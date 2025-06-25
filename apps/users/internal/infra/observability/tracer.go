package observability

import (
	"context"
	"fmt"

	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

// NewTracer initializes the OpenTelemetry TracerProvider.
// It sets up an OTLP gRPC exporter to send traces to a collector.
func NewTracer(ctx context.Context, config *config.Config) (*sdktrace.TracerProvider, error) {
	// Exporter configuration (e.g., to an OpenTelemetry Collector)
	traceExporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithInsecure(), // For local/development, use WithTLSCredentials for production
		otlptracegrpc.WithEndpoint("localhost:4317"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}

	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL, // ใช้ SchemaURL จาก semconv ที่ Import มา
			semconv.ServiceName(config.Name),
			semconv.ServiceVersion(config.Version),
			attribute.String("environment", "development"), // Or get from env variable
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()), // Sample all traces for simplicity in dev
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(sdktrace.NewBatchSpanProcessor(traceExporter)),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	return tp, nil
}
