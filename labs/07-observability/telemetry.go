package observability

import (
	"context"
	"log/slog"

	"go.opentelemetry.io/otel/trace"
)

// LogAttrsFromContext extracts trace_id and span_id from the active OpenTelemetry span.
func LogAttrsFromContext(ctx context.Context) []slog.Attr {
	span := trace.SpanFromContext(ctx)
	if !span.SpanContext().IsValid() {
		return nil
	}
	return []slog.Attr{
		slog.String("trace_id", span.SpanContext().TraceID().String()),
		slog.String("span_id", span.SpanContext().SpanID().String()),
	}
}

// ContextLogger returns a slog.Logger annotated with request_id, trace_id, and span_id.
func ContextLogger(ctx context.Context, base *slog.Logger, requestID string) *slog.Logger {
	if base == nil {
		base = slog.Default()
	}
	args := make([]any, 0, 6)
	if requestID != "" {
		args = append(args, "request_id", requestID)
	}
	for _, attr := range LogAttrsFromContext(ctx) {
		args = append(args, attr.Key, attr.Value.Any())
	}
	return base.With(args...)
}
