package main

import "errors"

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
	Key    string
	Value  any `json:"value"`
	Expiry int `json:"exp,omitempty"`
}

type PatchKeyRequest struct {
	Key    string
	Value  any `json:"value,omitempty"`
	Expiry int `json:"exp,omitempty"`
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

func (body PostKeyRequest) Validate() error {
	errs := make(map[string]string)

	if body.Key == "" {
		panic(errors.New("key cannot be empty"))
	}

	if body.Value == nil {
		errs["value"] = "value cannot be nil or empty"
	}

	if body.Expiry < 0 {
		errs["exp"] = "exp cannot be less than zero"
	}

	if len(errs) > 0 {
		return &BadRequestError{Message: "invalid request body", Fields: errs}
	}

	return nil
}

func (body PatchKeyRequest) Validate() error {
	errs := make(map[string]string)

	if body.Key == "" {
		panic(errors.New("key cannot be empty"))
	}

	if body.Value == nil && body.Expiry < 0 {
		errs["value"] = "value cannot be nil or empty"
		errs["exp"] = "exp cannot be less than zero"
	} else if body.Expiry < 0 {
		errs["exp"] = "exp cannot be less than zero"
	}

	if len(errs) > 0 {
		return &BadRequestError{Message: "invalid request body", Fields: errs}
	}

	return nil
}

func (body GetKeyRequest) Validate() error {
	if body.Key == "" {
		panic(errors.New("key cannot be empty"))
	}

	return nil
}

func (body DeleteKeyRequest) Validate() error {
	if body.Key == "" {
		panic(errors.New("key cannot be empty"))
	}

	return nil
}

func (body KeyExistsRequest) Validate() error {
	if body.Key == "" {
		panic(errors.New("key cannot be empty"))
	}

	return nil
}
