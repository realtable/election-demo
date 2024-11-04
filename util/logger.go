package util

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/middleware"
	"go.opentelemetry.io/otel/trace"
)

// TODO comments (incl explaining how to log and when to include context)

type TracedJSONHandler struct {
	slog.JSONHandler
}

func NewTracedJSONHandler(w *os.File, options *slog.HandlerOptions) *TracedJSONHandler {
	return &TracedJSONHandler{
		JSONHandler: *slog.NewJSONHandler(w, options),
	}
}

func (h *TracedJSONHandler) Handle(ctx context.Context, r slog.Record) error {
	spanContext := trace.SpanContextFromContext(ctx)
	if spanContext.HasTraceID() {
		r.AddAttrs(
			slog.String("traceID", spanContext.TraceID().String()),
			slog.String("spanID", spanContext.SpanID().String()),
		)
	}
	return h.JSONHandler.Handle(ctx, r)
}

func SetDefaultSlogger() {
	slogger := slog.New(NewTracedJSONHandler(os.Stdout, nil))
	slog.SetDefault(slogger)
}

type RequestLogFormatter struct{}
type requestLogEntry struct {
	r *http.Request
}

func (l *RequestLogFormatter) NewLogEntry(r *http.Request) middleware.LogEntry {
	return &requestLogEntry{
		r: r,
	}
}

func (l *requestLogEntry) Write(status, bytes int, header http.Header, elapsed time.Duration, extra interface{}) {
	slog.InfoContext(
		l.r.Context(),
		"request completed",
		"status", status,
		"bytes", bytes,
		"elapsedNs", elapsed,
		"path", l.r.URL.Path,
		"method", l.r.Method,
	)
}

func (l *requestLogEntry) Panic(v interface{}, stack []byte) {
	slog.Error("request panicked", "error", v, "stack", string(stack))
}
