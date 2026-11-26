package main

import (
	"context"
	"crypto/tls"
	"encoding/binary"
	"io"

	"github.com/google/uuid"
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
	RequestIdBytesSize           int = 16
)

const (
	CurrentVersion      uint8  = 1
	MaxPayloadSizeBytes uint32 = 17 * 1024 * 1024
	MaxKeySizeBytes     uint32 = 64 * 1024
	MaxValueSizeBytes   uint32 = 16 * 1024 * 1024
)

//go:generate stringer -type=Command
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
	RequestId   string
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
	RequestID   string
	Command     Command
	PayloadSize uint16
	KeyLength   uint32
	Key         string
	ValueLength uint32
	ValueType   byte
	ValueBytes  []byte
	Expiry      int
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
	InvalidHeadersOCPError
	InvalidPayloadOCPError
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
		ValueBytes:  ocpRequestBody.Value,
	}
}

func parseHeaderBytes(data []byte) (*OCPRequestHeader, *OCPError) {
	index := 0

	versionBytes := data[index : index+VersionHeaderLengthBytes]
	versionByte := versionBytes[0]
	index = index + VersionHeaderLengthBytes

	requestIdBytes := data[index : index+RequestIdBytesSize]
	index = index + RequestIdBytesSize
	requestId, err := uuid.FromBytes(requestIdBytes)
	if err != nil {
		return &OCPRequestHeader{}, &OCPError{ocpErrorType: InvalidHeadersOCPError}
	}
	parsedRequestId := requestId.String()

	commandBytes := data[index : index+CommandHeaderLengthBytes]
	commandByte := commandBytes[0]
	index = index + CommandHeaderLengthBytes

	payloadBytes := data[index : index+PayloadSizeHeaderLengthBytes]

	version := uint8(versionByte)
	if version != CurrentVersion {
		return &OCPRequestHeader{}, &OCPError{ocpErrorType: InvalidHeadersOCPError}
	}

	parsedCommand := int8(commandByte)
	if parsedCommand < 1 || parsedCommand > 5 {
		return &OCPRequestHeader{}, &OCPError{ocpErrorType: InvalidHeadersOCPError}
	}
	command := Command(parsedCommand)

	payloadSize := binary.BigEndian.Uint16(payloadBytes)
	if uint32(payloadSize) > MaxPayloadSizeBytes {
		return &OCPRequestHeader{}, &OCPError{ocpErrorType: InvalidHeadersOCPError}
	}

	return &OCPRequestHeader{
		Version:     version,
		RequestId:   parsedRequestId,
		Command:     command,
		PayloadSize: payloadSize,
	}, nil
}

func parsePayloadBytes(data []byte) (*OCPRequestBody, *OCPError) {
	index := 0
	keyLengthBytes := data[index : index+KeySizeHeaderLengthBytes]
	keySize := uint32(binary.BigEndian.Uint16(keyLengthBytes))

	if keySize > MaxKeySizeBytes {
		return &OCPRequestBody{}, &OCPError{ocpErrorType: InvalidPayloadOCPError}
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
		return &OCPRequestBody{}, &OCPError{ocpErrorType: InvalidPayloadOCPError}
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

func loadCertificate(certFile string, keyFile string) (tls.Certificate, error) {
	return tls.LoadX509KeyPair(certFile, keyFile)
}

const (
	certFilePath = ""
	keyFilePath  = ""
)

func InitializeQuicListener(port string, config *quic.Config) *quic.Listener {
	certificate, error := loadCertificate(certFilePath, keyFilePath)
	if error != nil {
		panic(error)
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{
			certificate,
		},
		MinVersion: tls.VersionTLS13,
		NextProtos: []string{"opc"},
	}

	listener, err := quic.ListenAddr(
		port,
		tlsConfig,
		config,
	)
	if err != nil {
		panic(err)
	}

	return listener
}

func ListenAndServe(listener *quic.Listener, handler *OCPHandler, middlewares []Middleware) error {
	for {
		conn, err := listener.Accept(context.Background())

		go handleConnection(conn, handler, middlewares)
		return err
	}
}

func handleConnection(conn *quic.Conn, handler *OCPHandler, middlewares []Middleware) {
	defer conn.CloseWithError(0, "connection closed")

	for {
		stream, err := conn.AcceptStream(context.Background())
		if err != nil {
			return
		}

		go handleStream(stream, handler, middlewares)
	}
}

func handleStream(stream *quic.Stream, handler *OCPHandler, middlewares []Middleware) {
	defer stream.Close()

	headerBuffer := make([]byte, TotalHeaderLengthBytes)

	_, err := io.ReadFull(stream, headerBuffer)
	if err != nil {
		// Handle incomplete header or stream error
		return
	}
	parsedRequestHeaders, headerErr := parseHeaderBytes(headerBuffer)
	if headerErr != nil {
		writeOCPErrorResponse(stream, uint8(headerErr.ocpErrorType))
		return
	}

	payloadBuffer := make([]byte, parsedRequestHeaders.PayloadSize)
	_, err = io.ReadFull(stream, payloadBuffer)
	if err != nil {
		return
	}

	parsedRequestPayload, payloadErr := parsePayloadBytes(payloadBuffer)
	if err != nil {
		writeOCPErrorResponse(stream, uint8(payloadErr.ocpErrorType))
		return
	}

	request := constructOCPRequest(parsedRequestHeaders, parsedRequestPayload)

	ocpResponse, ocpErr := runMiddlwares(middlewares, stream, handler, request)
	if ocpErr != nil {
		writeOCPErrorResponse(stream, uint8(ocpErr.ocpErrorType))
		return
	}

	writeOCPResponse(stream, 0, ocpResponse)
}
