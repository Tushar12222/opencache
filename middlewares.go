package main

import (
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/quic-go/quic-go"
)

type Middleware func(*quic.Stream, *OCPHandler, *OCPRequest) (*OCPResponseBody, *OCPError)

func extractRequestID(r *OCPRequest) string {
	requestID := r.RequestID
	return requestID
}

func validateRequestID(requestID string) *OCPError {
	_, err := uuid.Parse(requestID)
	if err != nil {
		return &OCPError{ocpErrorType: InvalidHeadersOCPError}
	}

	return nil
}

func LoggingMiddleware(stream *quic.Stream, handler OCPHandler, req *OCPRequest) (*OCPResponseBody, *OCPError) {
	requestID := extractRequestID(req)
	command := req.Command
	start := time.Now()

	reqLogger := slog.Default().With(
		slog.String("request_id", requestID),
		slog.String("command", command.String()),
	)

	reqLogger.Info("stream request received")

	response, err := handler.handleOCPRequest(stream, req)

	reqLogger.Info("stream response completed",
		slog.Duration("latency", time.Since(start)),
	)

	return response, err
}

func RequestIDMiddleware(stream *quic.Stream, handler OCPHandler, req *OCPRequest) (*OCPResponseBody, *OCPError) {
	requestID := extractRequestID(req)
	err := validateRequestID(requestID)
	if err != nil {
		return &OCPResponseBody{}, &OCPError{ocpErrorType: InvalidHeadersOCPError}
	}

	return handler.handleOCPRequest(stream, req)
}

func PanicRecoveryMiddleware(stream *quic.Stream, handler OCPHandler, req *OCPRequest) (*OCPResponseBody, *OCPError) {
	defer func() {
		if err := recover(); err != nil {
			slog.Error("PANIC RECOVERED", slog.Any("error", err))
			writeOCPErrorResponse(stream, uint8(InternalServerOCPError))
		}
	}()
	return handler.handleOCPRequest(stream, req)
}

func runMiddlwares(middlewares []Middleware, stream *quic.Stream, handler *OCPHandler, request *OCPRequest) (*OCPResponseBody, *OCPError) {
	var response *OCPResponseBody
	var err *OCPError
	for _, middleware := range middlewares {
		response, err = middleware(stream, handler, request)
	}

	return response, err
}
