package trace

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

func TraceIdFromContext(ctx context.Context) string {
	return trace.SpanFromContext(ctx).SpanContext().TraceID().String()
}
