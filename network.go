package main

import (
	"context"
	"crypto/tls"
	"encoding/binary"
	"errors"
	"fmt"
	"io"

	"github.com/quic-go/quic-go"
)

var (
	versionHeaderLengthBytes     int   = 1
	payloadSizeHeaderLengthBytes int   = 3
	commandHeaderLengthBytes     int   = 1
	totalHeaderLengthBytes       int   = versionHeaderLengthBytes + payloadSizeHeaderLengthBytes + commandHeaderLengthBytes
	currentVersion               uint8 = 1
)

type Command uint8

const (
	GetCommand Command = iota + 1
	PutCommant
	PostCommand
	DeleteCommand
	ExistsCommand
)

type OCPRequestHeader struct {
	Version     uint8
	Command     Command
	PayloadSize uint16
}

type OCPRequestPayload struct {
	KeyLength   int
	Key         string
	ValueLength int
	Value       any
}

func parseHeaderBytes(data []byte) (*OCPRequestHeader, error) {
	index := 0

	versionBytes := data[index:versionHeaderLengthBytes]
	versionByte := versionBytes[0]
	index = index + versionHeaderLengthBytes

	commandBytes := data[index:commandHeaderLengthBytes]
	commandByte := commandBytes[0]
	index = index + commandHeaderLengthBytes

	payloadBytes := data[index:payloadSizeHeaderLengthBytes]

	version := uint8(versionByte)
	if version != currentVersion {
		return &OCPRequestHeader{}, errors.New("")
	}

	parsedCommand := int8(commandByte)
	if parsedCommand < 1 || parsedCommand > 5 {
		return &OCPRequestHeader{}, errors.New("")
	}

	command := Command(parsedCommand)
	payloadSize := binary.BigEndian.Uint16(payloadBytes)

	return &OCPRequestHeader{
		Version:     version,
		Command:     command,
		PayloadSize: payloadSize,
	}, nil
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
	parsedRequestHeaders, err := parseHeaderBytes(headerBuffer)
	if err != nil {
		return
	}

	payloadBuffer := make([]byte, parsedRequestHeaders.PayloadSize)
	_, err = io.ReadFull(stream, payloadBuffer)
	if err != nil {
		return
	}

	parsedRequestPayload := parsePayloadBytes(payloadBuffer)

	fmt.Println("Received a new OPC stream")
}
