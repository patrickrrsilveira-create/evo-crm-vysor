package telemetry

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.27.0"
)

// InitTracer initializes OpenTelemetry tracing and returns a shutdown function
func InitTracer(ctx context.Context) (func(context.Context) error, error) {
	// Check if tracing is enabled
	if os.Getenv("OTEL_TRACES_ENABLED") != "true" {
		log.Println("OpenTelemetry tracing is disabled")
		return func(context.Context) error { return nil }, nil
	}

	// Get service name from environment
	serviceName := os.Getenv("OTEL_SERVICE_NAME")
	if serviceName == "" {
		serviceName = "evo-ai-core-service"
	}

	// Get OTLP endpoint
	otlpEndpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if otlpEndpoint == "" {
		return nil, fmt.Errorf("OTEL_EXPORTER_OTLP_ENDPOINT is not set")
	}

	// Parse endpoint URL to extract host:port
	// otlpEndpoint format: http://tempo.evo-monitoring:4318
	endpointURL, err := url.Parse(otlpEndpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to parse OTLP endpoint: %w", err)
	}

	// Create resource with service name
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(serviceName),
			attribute.String("deployment.environment", "production"),
		),
		resource.WithFromEnv(),
		resource.WithProcess(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// Create OTLP HTTP exporter
	exporter, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint(endpointURL.Host),
		otlptracehttp.WithURLPath("/v1/traces"),
		otlptracehttp.WithInsecure(), // Tempo uses HTTP, not HTTPS
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP exporter: %w", err)
	}

	// Create tracer provider
	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(res),
	)

	// Set global tracer provider
	otel.SetTracerProvider(tp)

	// Set global propagator
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	log.Printf("OpenTelemetry tracing initialized for service: %s, endpoint: %s", serviceName, otlpEndpoint)

	// Return shutdown function
	return tp.Shutdown, nil
}
