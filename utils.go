package main

import (
	"encoding/json"
	"errors"
	"io"
)

type expectedRequestBody interface {
	PostKeyRequest | PatchKeyRequest
}

func JsonDecoder[requestBody expectedRequestBody](body *io.ReadCloser) (*requestBody, error) {
	var dto requestBody
	decoder := json.NewDecoder(*body)
	if err := decoder.Decode(&dto); err != nil {
		return nil, errors.New("Unknown body")
	}

	return &dto, nil
}
