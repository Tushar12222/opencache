package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/google/uuid"
)

type expectedRequestBody interface {
	PostKeyRequest | PatchKeyRequest
}

type ResponseBody interface {
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

func FormatError(err any) error {
	var e error
	switch concreteErr := err.(type) {
	case *NotFoundError:
		e = concreteErr

	case *ConflictError:
		e = concreteErr

	case *BadRequestError:
		e = concreteErr

	case *InternalServerError:
		concreteErr.Message = "something went wrong"

	default:
		e = &InternalServerError{Message: "something went wrong"}
	}
	return e
}

func writeErrorResponse(w http.ResponseWriter, status int, errMsg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err := w.Write([]byte(errMsg))
	if err != nil {
		panic(err)
	}
}

func writeResponse[responseBody ResponseBody](w http.ResponseWriter, status int, resBody responseBody) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(resBody)
	if err != nil {
		panic(err)
	}
}

func mustGenerateUUIDV7(
	generate func() (uuid.UUID, error),
) string {
	uuid, err := generate()
	if err != nil {
		panic(err)
	}
	return uuid.String()
}

func JSONEncoder(data any) ([]byte, error) {
	jsonData, err := json.Marshal(data)
	return jsonData, err
}
