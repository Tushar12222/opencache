package main

import (
	"context"
	"log/slog"
	"os"
)

type ctxKey string

const (
	KeyTraceID   ctxKey = "trace_id"
	KeyRequestID ctxKey = "request_id"
)

type ContextHandler struct {
	slog.Handler
}

func newJSONHandler() slog.Handler {
	return slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: true,
	})
}

func (h *ContextHandler) Handle(ctx context.Context, r slog.Record) error {
	if traceID, ok := ctx.Value(KeyTraceID).(string); ok {
		r.AddAttrs(slog.String("trace_id", traceID))
	}
	if reqID, ok := ctx.Value(KeyRequestID).(string); ok {
		r.AddAttrs(slog.String("request_id", reqID))
	}
	return h.Handler.Handle(ctx, r)
}

func InitProductionLogger() {
	baseJSONHandler := newJSONHandler()
	wrappedHandler := &ContextHandler{Handler: baseJSONHandler}

	slog.SetDefault(slog.New(wrappedHandler))
}

func NewCronLogger(jobName string) *slog.Logger {
	return slog.New(newJSONHandler()).With(
		slog.String("component", "cron"),
		slog.String("job", jobName),
	)
}
