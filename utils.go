package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type expectedRequestBody interface {
	PostKeyRequest | PatchKeyRequest
}

func JSONDecoder[requestBody expectedRequestBody](body *io.ReadCloser) (*requestBody, error) {
	var dto requestBody
	decoder := json.NewDecoder(*body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&dto); err != nil {
		return nil, errors.New("unknown body")
	}

	return &dto, nil
}

func writeErrorResponse(w http.ResponseWriter, err error) {
	errMsg := err.Error()
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(errMsg))
}

func writeResponse(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Received Request"))
}
