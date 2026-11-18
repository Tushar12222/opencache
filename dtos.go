package main

import (
	"errors"
	"time"
)

type GetKeyRequest struct {
	Key string
}

type KeyExistsRequest struct {
	Key string
}

type DeleteKeyRequest struct {
	Key string
}

type PostKeyRequest struct {
	Key        string
	ValueType  byte
	ValueBytes []byte
	Expiry     int
}

type PatchKeyRequest struct {
	Key        string
	ValueType  byte
	ValueBytes []byte
	Expiry     int
}

type PostKeyResponse struct {
	Message any `json:"message"`
}

type PatchKeyResponse struct {
	Message any `json:"message"`
}

type GetKeyResponse struct {
	Message any            `json:"message"`
	Data    map[string]any `json:"data"`
}

type DeleteKeyResponse struct {
	Message any `json:"message"`
}

type KeyExistsResponse struct {
	Message any            `json:"message"`
	Data    map[string]any `json:"data"`
}

type OCPResponseBody struct {
	ValueType   byte
	ValueLength uint32
	ValueBytes  []byte
	Expiration  *time.Time
}

func toOCPResponseBody(valueType byte, valueLength uint32, valueBytes []byte, expiration *time.Time) (*OCPResponseBody, *OCPError) {
	l := uint32(len(valueBytes))

	if l > valueLength {
		return &OCPResponseBody{}, &OCPError{InternalServerOCPError}
	}

	return &OCPResponseBody{
		valueType,
		valueLength,
		valueBytes,
		expiration,
	}, nil
}

func (body PostKeyRequest) Validate() *OCPError {
	errs := make(map[string]string)

	if body.Key == "" {
		panic(errors.New("key cannot be empty"))
	}

	if body.ValueType == 0x00 {
		errs["valueType"] = "valueType cannot be 0"
	}

	if len(body.ValueBytes) == 0 || body.ValueBytes == nil {
		errs["valueBytes"] = "valueBytes cannot be empty or nil"
	}

	if body.Expiry < 0 {
		errs["exp"] = "exp cannot be less than zero"
	}

	if len(errs) > 0 {
		return &OCPError{ocpErrorType: BadRequestOCPError}
	}

	return nil
}

func (body PatchKeyRequest) Validate() *OCPError {
	errs := make(map[string]string)

	if body.Key == "" {
		panic(errors.New("key cannot be empty"))
	}

	if (body.ValueBytes == nil || len(body.ValueBytes) == 0) && body.Expiry < 0 {
		errs["value"] = "value cannot be nil or empty"
		errs["exp"] = "exp cannot be less than zero"
	} else if body.Expiry < 0 {
		errs["exp"] = "exp cannot be less than zero"
	}

	if len(errs) > 0 {
		return &OCPError{ocpErrorType: BadRequestOCPError}
	}

	return nil
}

func (body GetKeyRequest) Validate() *OCPError {
	if body.Key == "" {
		panic(errors.New("key cannot be empty"))
	}

	return nil
}

func (body DeleteKeyRequest) Validate() *OCPError {
	if body.Key == "" {
		panic(errors.New("key cannot be empty"))
	}

	return nil
}

func (body KeyExistsRequest) Validate() *OCPError {
	if body.Key == "" {
		panic(errors.New("key cannot be empty"))
	}

	return nil
}
