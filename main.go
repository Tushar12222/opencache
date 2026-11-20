package main

import (
	"log/slog"
	"strconv"
)

func main() {
	stop := InitializeGlobalContext()
	defer stop()

	handler := InitProductionLogger()
	defer handler.Shutdown()

	quicListener := InitializeQuicListener(":"+strconv.Itoa(ServerPort), nil)
	quicHandler := newOCPHandler()
	StartCronJob(quicHandler.Service.cache, GetGlobalContext(), removeDeadKeysThread)

	serverErr := make(chan error, 1)

	go func() {
		slog.Info("Server running",
			slog.Int("port_number", ServerPort))
		serverErr <- ListenAndServe(quicListener, quicHandler)
	}()

	ListenForInterrupt(serverErr)
	cancel := CleanUp(quicListener)
	defer cancel()
	slog.Info("Server shutdown successfully")
}
