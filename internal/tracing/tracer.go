package tracing

import (
	"context"

	"github.com/webitel/cases/model"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var tracer trace.Tracer

// InitTracer initializes the global tracer.
func InitTracer() {
	tp := otel.GetTracerProvider()
	tracer = tp.Tracer(model.APP_SERVICE_NAME)
}

// StartSpan starts a new span with the provided context and operation name.
func StartSpan(ctx context.Context, operationName string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	return tracer.Start(ctx, operationName, opts...)
}

// GetTracer returns the initialized tracer.
func GetTracer() trace.Tracer {
	return tracer
}
