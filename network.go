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

const (
	VersionHeaderLengthBytes     int = 1
	PayloadSizeHeaderLengthBytes int = 3
	CommandHeaderLengthBytes     int = 1
	TotalHeaderLengthBytes       int = VersionHeaderLengthBytes + PayloadSizeHeaderLengthBytes + CommandHeaderLengthBytes
	KeySizeHeaderLengthBytes     int = 3
	ValueSizeHeaderLengthBytes   int = 4
	ValueTypeHeaderLengthBytes   int = 1
)

const (
	CurrentVersion      uint8  = 1
	MaxPayloadSizeBytes uint32 = 17 * 1024 * 1024
	MaxKeySizeBytes     uint32 = 64 * 1024
	MaxValueSizeBytes   uint32 = 16 * 1024 * 1024
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

type OCPRequestBody struct {
	KeyLength   uint32
	Key         string
	ValueLength uint32
	ValueType   byte
	Value       []byte
}

type OCPPayload struct {
	Version     uint8
	Command     Command
	PayloadSize uint16
	KeyLength   uint32
	Key         string
	ValueLength uint32
	ValueType   byte
	Value       []byte
}

// Network response for OCP
/*
[Version:      1 byte]
[Status:       1 byte]
[ResponseType: 1 byte]
[ValueType:    1 byte]
[ValueLength:  4 bytes]
[Value:        N bytes]

[Version:      1 byte]
[Status:       1 byte]
[ErrorCode: 1 byte]
[MessageLength: 3 bytes]
[Message: N bytes]
* */

func constructOCPPayload(ocpRequestHeader *OCPRequestHeader, ocpRequestBody *OCPRequestBody) *OCPPayload {
	return &OCPPayload{
		Version:     ocpRequestHeader.Version,
		Command:     ocpRequestHeader.Command,
		PayloadSize: ocpRequestHeader.PayloadSize,
		KeyLength:   ocpRequestBody.KeyLength,
		Key:         ocpRequestBody.Key,
		ValueLength: ocpRequestBody.ValueLength,
		ValueType:   ocpRequestBody.ValueType,
		Value:       ocpRequestBody.Value,
	}
}

func parseHeaderBytes(data []byte) (*OCPRequestHeader, error) {
	index := 0

	versionBytes := data[index : index+VersionHeaderLengthBytes]
	versionByte := versionBytes[0]
	index = index + VersionHeaderLengthBytes

	commandBytes := data[index : index+CommandHeaderLengthBytes]
	commandByte := commandBytes[0]
	index = index + CommandHeaderLengthBytes

	payloadBytes := data[index : index+PayloadSizeHeaderLengthBytes]

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
	if uint32(payloadSize) > MaxPayloadSizeBytes {
		return &OCPRequestHeader{}, errors.New("")
	}

	return &OCPRequestHeader{
		Version:     version,
		Command:     command,
		PayloadSize: payloadSize,
	}, nil
}

func parsePayloadBytes(data []byte) (*OCPRequestBody, error) {
	index := 0
	keyLengthBytes := data[index : index+KeySizeHeaderLengthBytes]
	keySize := uint32(binary.BigEndian.Uint16(keyLengthBytes))

	if keySize > MaxKeySizeBytes {
		return &OCPRequestBody{}, errors.New("")
	}
	index += KeySizeHeaderLengthBytes

	keyBytes := data[index : index+int(keySize)]
	key := string(keyBytes)

	valueLengthBytes := data[index : index+ValueSizeHeaderLengthBytes]
	index += ValueSizeHeaderLengthBytes

	valueTypeBytes := data[index : index+ValueTypeHeaderLengthBytes]
	index += ValueTypeHeaderLengthBytes
	valueTypeByte := valueTypeBytes[0]

	valueSize := binary.BigEndian.Uint32(valueLengthBytes)
	if valueSize > MaxValueSizeBytes {
		return &OCPRequestBody{}, errors.New("")
	}

	valueBytes := data[index : index+int(valueSize)]

	return &OCPRequestBody{
		KeyLength:   keySize,
		Key:         key,
		ValueLength: valueSize,
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

	headerBuffer := make([]byte, TotalHeaderLengthBytes)

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
