package main

import (
	"context"
	"crypto/tls"
	"encoding/binary"
	"errors"
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
	PatchCommand
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

type OCPRequest struct {
	Version     uint8
	Command     Command
	PayloadSize uint16
	KeyLength   uint32
	Key         string
	ValueLength uint32
	ValueType   byte
	Value       []byte
}

type OCPResponse struct {
	Version       uint8
	Status        uint8
	ResponseType  uint8
	ValueType     byte
	ValueLength   uint32
	Value         []byte
	ErrorCode     uint8
	MessageLength uint32
	Message       []byte
}

type OCPError struct {
	ocpErrorType OCPErrorType
}

type OCPErrorType uint8

const (
	KeyNotFoundOCPError OCPErrorType = iota + 1
	DuplicateKeyOCPError
	BadRequestOCPError
	InternalServerOCPError
)

func constructOCPRequest(ocpRequestHeader *OCPRequestHeader, ocpRequestBody *OCPRequestBody) *OCPRequest {
	return &OCPRequest{
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

func StartQuicServer() {
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

	handler := newOCPHandler()

	for {
		conn, err := listener.Accept(context.Background())
		if err != nil {
			panic(err)
		}

		go handleConnection(conn, handler)
	}
}

func handleConnection(conn *quic.Conn, handler *OCPHandler) {
	defer conn.CloseWithError(0, "connection closed")

	for {
		stream, err := conn.AcceptStream(context.Background())
		if err != nil {
			return
		}

		go handleStream(stream, handler)
	}
}

func (handler *OCPHandler) handleOCPRequest(w *quic.Stream, request *OCPRequest) {
	switch request.Command {
	case PatchCommand:
		handler.handleOCPPatchKey(w, request)
	case GetCommand:
		handler.handleOCPGetKey(w, request)
	case PostCommand:
		handler.handleOCPPostKey(w, request)
	case DeleteCommand:
		handler.handleOCPDeleteKey(w, request)
	case ExistsCommand:
		handler.handleOCPKeyExists(w, request)
	}
}

func handleStream(stream *quic.Stream, handler *OCPHandler) {
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

	parsedRequestPayload, err := parsePayloadBytes(payloadBuffer)
	if err != nil {
		return
	}

	request := constructOCPRequest(parsedRequestHeaders, parsedRequestPayload)
	handler.handleOCPRequest(stream, request)
}
