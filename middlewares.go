package main

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/quic-go/quic-go"
)

type Middleware func(*quic.Stream, *OCPHandler, *OCPRequest) (*OCPResponseBody, *OCPError)

type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriterWrapper) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func extractRequestID(r *http.Request) string {
	requestID := r.Header.Get("X-Request-ID")
	return requestID
}

func validateRequestID(requestID string) error {
	_, err := uuid.Parse(requestID)
	if err != nil {
		return &BadRequestError{Message: "invalid request", Fields: map[string]string{
			"requestID": err.Error(),
		}}
	}

	return nil
}

func LoggingMiddleware(stream *quic.Stream, handler OCPHandler, req *OCPRequest) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := extractRequestID(r)
		start := time.Now()

		wrappedWriter := &responseWriterWrapper{ResponseWriter: w, statusCode: http.StatusOK}

		reqLogger := slog.Default().With(
			slog.String("request_id", requestID),
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
		)

		reqLogger.Info("http request received")

		handler.handleOCPRequest(stream, req)

		reqLogger.Info("http request completed",
			slog.Int("status", wrappedWriter.statusCode),
			slog.Duration("latency", time.Since(start)),
		)
	})
}

func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := extractRequestID(r)
		if requestID == "" {
			requestID = mustGenerateUUIDV7(uuid.NewV7)
		} else {
			err := validateRequestID(requestID)
			if err != nil {
				writeErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
		}
		r.Header.Set("X-Request-ID", requestID)
		next.ServeHTTP(w, r)
	})
}

func PanicRecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				slog.Error("PANIC RECOVERED", slog.Any("error", err))
				errMsg, _ := JSONEncoder(map[string]any{
					"message": "something went wrong",
				})
				writeErrorResponse(w, http.StatusInternalServerError, string(errMsg))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func runMiddlwares(middlewares []Middleware, stream *quic.Stream, handler *OCPHandler, request *OCPRequest) (*OCPResponseBody, *OCPError) {
	var response *OCPResponseBody
	var err *OCPError
	for _, middleware := range middlewares {
		response, err = middleware(stream, handler, request)
	}

	return response, err
}
