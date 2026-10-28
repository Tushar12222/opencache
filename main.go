package main

import (
	"log/slog"
)

func main() {
	stop := InitializeGlobalContext()
	defer stop()

	handler := InitProductionLogger()
	defer handler.Shutdown()

	StartCronJob(GetGlobalHandler(), GetGlobalContext())
	httpServer := InitializeHTTPServer()

	serverErr := make(chan error, 1)

	go func() {
		slog.Info("Server running",
			slog.Int("port_number", ServerPort))
		serverErr <- httpServer.ListenAndServe()
	}()

	ListenForInterrupt(httpServer, serverErr)
	cancel := CleanUp(httpServer)
	defer cancel()
	slog.Info("Server shutdown successfully")
}
