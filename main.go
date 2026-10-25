package main

import (
	"fmt"
)

func main() {
	stop := InitializeGlobalContext()
	defer stop()

	InitProductionLogger()

	httpServer := InitializeHTTPServer()

	serverErr := make(chan error, 1)

	go func() {
		fmt.Println("Server running on port", ServerPort)
		serverErr <- httpServer.ListenAndServe()
	}()

	ListenForInterrupt(httpServer, serverErr)
	cancel := CleanUp(httpServer)
	defer cancel()
	fmt.Println("Server shutdown successfully")
}
