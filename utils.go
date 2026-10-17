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

type expectedResponseBody interface {
	PostKeyResponse | PatchKeyResponse | GetKeyResponse | DeleteKeyResponse
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

func writeErrorResponse(w http.ResponseWriter, status int, errMsg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{
		"message": errMsg,
	})
}

func writeResponse[responseBody expectedResponseBody](w http.ResponseWriter, status int, resBody responseBody) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(resBody)
}
