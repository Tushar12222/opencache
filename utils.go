package main

import (
	"encoding/json"
	"errors"
	"net/http"
)

type expectedRequestBody interface {
	PostKeyRequest | PatchKeyRequest
}

func JsonDecoder[requestBody expectedRequestBody](decoder *json.Decoder, request *http.Request) (*requestBody, error) {
	var dto requestBody
	if err := decoder.Decode(&dto); err != nil {
		return nil, errors.New("Unknown body")
	}

	return &dto, nil
}
