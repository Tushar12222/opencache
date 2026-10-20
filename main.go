package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

var (
	ServerPort   = 8080
	ReadTimeout  = 5 * time.Second
	WriteTimeout = 10 * time.Second
	IdleTimeout  = 120 * time.Second
)

func main() {
	mux := http.NewServeMux()

	RegisterRoutes(mux)

	globalWrappedMux := PanicRecoveryMiddleware(RequestIDMiddleware(LoggingMiddleware(mux)))

	server := &http.Server{
		Addr:         ":" + strconv.Itoa(ServerPort),
		Handler:      globalWrappedMux,
		ReadTimeout:  ReadTimeout,
		WriteTimeout: WriteTimeout,
		IdleTimeout:  IdleTimeout,
	}

	fmt.Println("Server running on port", ServerPort)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}
