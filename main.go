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
	// Create a dedicated route multiplexer (router)
	mux := http.NewServeMux()

	// Register routes
	RegisterRoutes(mux)

	// Explicitly configure server properties to prevent connection exhaustion
	server := &http.Server{
		Addr:         ":" + strconv.Itoa(ServerPort),
		Handler:      mux,
		ReadTimeout:  ReadTimeout,
		WriteTimeout: WriteTimeout,
		IdleTimeout:  IdleTimeout,
	}

	fmt.Println("Server running on port", ServerPort)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}
