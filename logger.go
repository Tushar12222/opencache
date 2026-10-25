package main

import (
	"context"
	"log/slog"
	"os"
)

type ctxKey string

const (
	KeyTraceID        ctxKey = "trace_id"
	KeyRequestID      ctxKey = "request_id"
	KeyIdempotencyKey ctxKey = "idempotency_key"
)

type ContextHandler struct {
	slog.Handler
}

func (h *ContextHandler) Handle(ctx context.Context, r slog.Record) error {
	if traceID, ok := ctx.Value(KeyTraceID).(string); ok {
		r.AddAttrs(slog.String("trace_id", traceID))
	}
	if reqID, ok := ctx.Value(KeyRequestID).(string); ok {
		r.AddAttrs(slog.String("request_id", reqID))
	}
	if idemKey, ok := ctx.Value(KeyIdempotencyKey).(string); ok {
		r.AddAttrs(slog.String("idempotency_key", idemKey))
	}
	return h.Handler.Handle(ctx, r)
}

func InitProductionLogger() {
	opts := &slog.HandlerOptions{
		Level:     slog.LevelInfo, // Suppress debug messages
		AddSource: true,           // Capture filename and line number
	}

	baseJSONHandler := slog.NewJSONHandler(os.Stdout, opts)
	wrappedHandler := &ContextHandler{Handler: baseJSONHandler}

	slog.SetDefault(slog.New(wrappedHandler))
}
