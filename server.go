package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
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

func shutDownHTTPServer(server *http.Server, ctx context.Context) {
	if err := server.Shutdown(ctx); err != nil {
		fmt.Println("Graceful shutdown failed:", err)
		if closeErr := server.Close(); closeErr != nil {
			fmt.Println("Force close failed:", closeErr)
		}
	}
}

func CleanUp(server *http.Server) context.CancelFunc {
	fmt.Println("Cleaning up resources")
	cleanUpTimeout := 10 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), cleanUpTimeout)
	shutDownHTTPServer(server, ctx)
	return cancel
}

func ListenForInterrupt(server *http.Server, serverErr chan error) {
	select {
	case <-GetGlobalContext().Done():
		fmt.Println("Interrupt received")

	case err := <-serverErr:
		if err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}
}

func StartBackgroundThread(handler *Handler, context context.Context) {
	go func() {
		removeDeadKeysThread(handler, context)
	}()
}

func removeDeadKeysThread(handler *Handler, context context.Context) {
	backgroundThreadInterval := 10 * time.Second
	ticker := time.NewTicker(backgroundThreadInterval)
	defer ticker.Stop()
	cache := handler.Service.cache
	for {
		select {
		case <-context.Done():
			return
		case <-ticker.C:
			cache.RemoveExpiredKeys()
		}
	}
}

func InitializeHTTPServer() *http.Server {
	mux := http.NewServeMux()

	globalWrappedMux := PanicRecoveryMiddleware(RequestIDMiddleware(LoggingMiddleware(mux)))

	RegisterRoutes(mux)

	StartBackgroundThread(GetGlobalHandler(), GetGlobalContext())

	server := &http.Server{
		Addr:         ":" + strconv.Itoa(ServerPort),
		Handler:      globalWrappedMux,
		ReadTimeout:  ReadTimeout,
		WriteTimeout: WriteTimeout,
		IdleTimeout:  IdleTimeout,
	}

	return server
}
