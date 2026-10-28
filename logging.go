package main

import (
	"context"
	"log/slog"
	"os"
	"sync"
)

type logEntry struct {
	ctx     context.Context
	handler slog.Handler
	record  slog.Record
}

type asyncState struct {
	queue chan logEntry
	mu    sync.Mutex
	wg    sync.WaitGroup
}

type AsyncHandler struct {
	inner slog.Handler
	state *asyncState
}

func NewAsyncHandler(inner slog.Handler, bufferSize int) *AsyncHandler {
	state := &asyncState{
		queue: make(chan logEntry, bufferSize),
	}

	state.wg.Add(1)
	go func() {
		defer state.wg.Done()

		for entry := range state.queue {
			_ = entry.handler.Handle(entry.ctx, entry.record)
		}
	}()

	return &AsyncHandler{
		inner: inner,
		state: state,
	}
}

func (h *AsyncHandler) Enabled(
	ctx context.Context,
	level slog.Level,
) bool {
	return h.inner.Enabled(ctx, level)
}

func (h *AsyncHandler) Handle(
	ctx context.Context,
	record slog.Record,
) error {
	h.state.mu.Lock()
	defer h.state.mu.Unlock()

	h.state.queue <- logEntry{
		ctx:     ctx,
		handler: h.inner,
		record:  record.Clone(),
	}

	return nil
}

func (h *AsyncHandler) WithAttrs(
	attrs []slog.Attr,
) slog.Handler {
	return &AsyncHandler{
		inner: h.inner.WithAttrs(attrs),
		state: h.state,
	}
}

func (h *AsyncHandler) WithGroup(name string) slog.Handler {
	return &AsyncHandler{
		inner: h.inner.WithGroup(name),
		state: h.state,
	}
}

func (h *AsyncHandler) Shutdown() {
	h.state.mu.Lock()
	close(h.state.queue)
	h.state.mu.Unlock()

	h.state.wg.Wait()
}

type ctxKey string

const (
	KeyTraceID      ctxKey = "trace_id"
	KeyRequestID    ctxKey = "request_id"
	InMemoryLogSize int    = 1024
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

func InitProductionLogger() *AsyncHandler {
	baseJSONHandler := newJSONHandler()
	wrappedHandler := &ContextHandler{Handler: baseJSONHandler}
	handler := NewAsyncHandler(wrappedHandler, InMemoryLogSize)
	slog.SetDefault(slog.New(handler))

	return handler
}

func NewCronLogger(jobName string, jobID string) *slog.Logger {
	return slog.New(newJSONHandler()).With(
		slog.String("component", "cron"),
		slog.String("job", jobName),
		slog.String("jobID", jobID),
	)
}
