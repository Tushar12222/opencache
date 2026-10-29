package main

import (
	"fmt"
	"io"
	"net"
)

func HandleAndServe(port string) (err error) {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		panic(err)
	}

	defer func() {
		if closeErr := listener.Close(); err != nil {
			err = closeErr
		}
	}()

	fmt.Println("TCP server listening on :8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		fmt.Println("Client connected:", conn.RemoteAddr())

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 1024)

	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				fmt.Println("Read error:", err)
			}
			return
		}

		fmt.Printf("Received: %s\n", buf[:n])

		_, err = conn.Write(buf[:n]) // Echo data back to client
		if err != nil {
			fmt.Println("Write error:", err)
			return
		}
	}
}
