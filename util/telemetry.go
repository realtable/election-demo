package util

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"strconv"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	oteltrace "go.opentelemetry.io/otel/trace"
)

var Tracer oteltrace.Tracer

var TelemetryMiddleware = otelhttp.NewMiddleware(
	"serveRoute", // default span name
	otelhttp.WithSpanNameFormatter(func(operation string, r *http.Request) string {
		return r.Method + " " + r.URL.Path
	}),
)

func InitTracerProvider() (*sdktrace.TracerProvider, error) {
	ctx := context.Background()
	serviceName := os.Getenv("SERVICE_NAME")

	exporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint(os.Getenv("OTLP_ENDPOINT")),
		otlptracegrpc.WithInsecure()) // TODO indent like this or with every arg on its own line?
	if err != nil {
		return nil, err
	}

	traceRatio, err := strconv.ParseFloat(os.Getenv("TRACE_RATIO"), 64)
	if err != nil {
		slog.Error("could not parse TRACE_RATIO", "error", err)
		os.Exit(1)
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.TraceIDRatioBased(traceRatio))),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String(serviceName),
			)),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	Tracer = otel.Tracer(serviceName)

	return tp, nil
}
