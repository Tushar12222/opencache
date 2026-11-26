package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/quic-go/quic-go"
)

var (
	ServerPort   = 8080
	ReadTimeout  = 5 * time.Second
	WriteTimeout = 10 * time.Second
	IdleTimeout  = 120 * time.Second
)

var globalContext context.Context

func InitializeGlobalContext() context.CancelFunc {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)

	globalContext = ctx
	return stop
}

func GetGlobalContext() context.Context {
	return globalContext
}

func shutDownQuicListener(listener *quic.Listener, ctx context.Context) {
	if err := listener.Close(); err != nil {
		fmt.Println("quic listener shutdown failed:", err)
	}
}

func CleanUp(listener *quic.Listener) context.CancelFunc {
	slog.Info("Cleaning up resources")
	cleanUpTimeout := 10 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), cleanUpTimeout)
	shutDownQuicListener(listener, ctx)
	return cancel
}

func ListenForInterrupt(serverErr chan error) {
	select {
	case <-GetGlobalContext().Done():
		slog.Info("Interrupt received")

	case err := <-serverErr:
		if err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}
}

func StartCronJob(cache *Cache, context context.Context, removeExpiredKeyCron func(*Cache, context.Context)) {
	go removeExpiredKeyCron(cache, context)
}
