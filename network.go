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
	versionHeaderLengthBytes     int = 1
	payloadSizeHeaderLengthBytes int = 3
	commandHeaderLengthBytes     int = 1
	totalHeaderLengthBytes       int = versionHeaderLengthBytes + payloadSizeHeaderLengthBytes + commandHeaderLengthBytes
	keySizeHeaderLengthBytes     int = 1
	valueSizeHeaderLengthBytes   int = 2
	valueTypeHeaderLengthBytes   int = 1
)

const (
	CurrentVersion uint8   = 1
	MaxPayloadSize uint16  = 17
	MaxKeySize     float32 = 0.64
	MaxValueSize   uint16  = 17
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
	ValueType   byte
	Value       []byte
}

func parseHeaderBytes(data []byte) (*OCPRequestHeader, error) {
	index := 0

	versionBytes := data[index : index+versionHeaderLengthBytes]
	versionByte := versionBytes[0]
	index = index + versionHeaderLengthBytes

	commandBytes := data[index : index+commandHeaderLengthBytes]
	commandByte := commandBytes[0]
	index = index + commandHeaderLengthBytes

	payloadBytes := data[index : index+payloadSizeHeaderLengthBytes]

	version := uint8(versionByte)
	if version != CurrentVersion {
		return &OCPRequestHeader{}, errors.New("")
	}

	parsedCommand := int8(commandByte)
	if parsedCommand < 1 || parsedCommand > 5 {
		return &OCPRequestHeader{}, errors.New("")
	}
	command := Command(parsedCommand)

	payloadSize := binary.BigEndian.Uint16(payloadBytes)
	if payloadSize > MaxPayloadSize {
		return &OCPRequestHeader{}, errors.New("")
	}

	return &OCPRequestHeader{
		Version:     version,
		Command:     command,
		PayloadSize: payloadSize,
	}, nil
}

func parsePayloadBytes(data []byte) (*OCPRequestPayload, error) {
	index := 0
	keyLengthBytes := data[index : index+keySizeHeaderLengthBytes]
	keySize := binary.BigEndian.Uint16(keyLengthBytes)

	if float32(keySize) > MaxKeySize {
		return &OCPRequestPayload{}, errors.New("")
	}
	index += keySizeHeaderLengthBytes

	keyBytes := data[index : index+int(keySize)]
	key := string(keyBytes)

	valueLengthBytes := data[index : index+valueSizeHeaderLengthBytes]
	index += valueSizeHeaderLengthBytes

	valueTypeBytes := data[index : index+valueTypeHeaderLengthBytes]
	index += valueTypeHeaderLengthBytes
	valueTypeByte := valueTypeBytes[0]

	valueSize := binary.BigEndian.Uint16(valueLengthBytes)
	if valueSize > MaxValueSize {
		return &OCPRequestPayload{}, errors.New("")
	}

	valueBytes := data[index : index+int(valueSize)]

	return &OCPRequestPayload{
		KeyLength:   int(keySize),
		Key:         key,
		ValueLength: int(valueSize),
		ValueType:   valueTypeByte,
		Value:       valueBytes,
	}, nil
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
