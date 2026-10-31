package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"

	"github.com/quic-go/quic-go"
)

var (
	versionHeaderLengthBytes     int = 1
	payloadSizeHeaderLengthBytes int = 1
	commandHeaderLengthBytes     int = 4
	totalHeaderLengthBytes       int = versionHeaderLengthBytes + payloadSizeHeaderLengthBytes + commandHeaderLengthBytes
)

type OCPRequestHeader struct {
	Version     string
	Command     string
	PayloadSize int
}

type OCPRequestPayload struct {
	KeyLength   int
	Key         string
	ValueLength int
	Value       any
}

func parseHeaderBytes(data []byte) *OCPRequestHeader {
	return nil
}

func parsePayloadBytes(data []byte) *OCPRequestPayload {
	return nil
}

func loadCertificate() {}

func StartQuicListener() {
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{
			loadCertificate(),
		},
		NextProtos: []string{"opc"},
	}

	listener, err := quic.ListenAddr(
		":4242",
		tlsConfig,
		nil,
	)
	if err != nil {
		panic(err)
	}

	for {
		conn, err := listener.Accept(context.Background())
		if err != nil {
			panic(err)
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn *quic.Conn) {
	defer conn.CloseWithError(0, "connection closed")

	for {
		stream, err := conn.AcceptStream(context.Background())
		if err != nil {
			return
		}

		go handleStream(stream)
	}
}

func handleStream(stream *quic.Stream) {
	defer stream.Close()

	headerBuffer := make([]byte, totalHeaderLengthBytes)

	_, err := io.ReadFull(stream, headerBuffer)
	if err != nil {
		// Handle incomplete header or stream error
		return
	}
	parsedRequestHeaders := parseHeaderBytes(headerBuffer)

	payloadBuffer := make([]byte, parsedRequestHeaders.PayloadSize)
	_, err = io.ReadFull(stream, payloadBuffer)
	if err != nil {
		return
	}

	parsedRequestPayload := parsePayloadBytes(payloadBuffer)

	fmt.Println("Received a new OPC stream")
}
