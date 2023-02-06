package tracing

import (
	"os"

	"github.com/stumacwastaken/todo/log"
	"go.opentelemetry.io/contrib/propagators/autoprop"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
)

var tracer trace.Tracer

func Tracer() trace.Tracer {
	return tracer
}
func newExporter() (tracesdk.SpanExporter, error) {
	exportURL := os.Getenv("OTEL_EXPORTER_JAEGER_ENDPOINT")

	log.Default().Info("setting tracing exporter to jaeger", zap.String("export-url", exportURL))
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(exportURL)))

	return exp, err
}

func newResource(serviceName string) *resource.Resource {
	return resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(serviceName),
	)
}

func InitTracingProvider(serviceName string) *tracesdk.TracerProvider {
	exp, err := newExporter()
	if err != nil {
		log.Default().Error("error creating tracer", zap.Error(err))
		return nil
	}
	tp := tracesdk.NewTracerProvider(
		tracesdk.WithSampler(tracesdk.AlwaysSample()),
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(newResource(serviceName)),
	)
	otel.SetTracerProvider(tp)
	tracer = tp.Tracer("todo")
	otel.SetTextMapPropagator(autoprop.NewTextMapPropagator())

	return tp
}
