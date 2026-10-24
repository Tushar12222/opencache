package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

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
		return nil, err
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

func parseDecodeError(err error) *BadRequestError {
	fields := make(map[string]string)

	var typeErr *json.UnmarshalTypeError
	var syntaxErr *json.SyntaxError
	switch {
	case errors.Is(err, io.EOF):
		fields["body"] = "request body is required"
	case errors.As(err, &typeErr):
		field := typeErr.Field
		if field == "" {
			field = "body"
		}
		fields[field] = fmt.Sprintf("expected %s, got %s", typeErr.Type.String(), typeErr.Value)
	case errors.As(err, &syntaxErr):
		fields["body"] = fmt.Sprintf("malformed JSON at byte %d", syntaxErr.Offset)
	default:
		if unknownField, ok := unknownJSONField(err.Error()); ok {
			fields[unknownField] = "unknown field"
		} else {
			fields["body"] = "invalid JSON request body"
		}
	}

	if len(fields) > 0 {
		return &BadRequestError{
			Message: "invalid request body",
			Fields:  fields,
		}
	}

	return nil
}

func unknownJSONField(errMsg string) (string, bool) {
	const prefix = "json: unknown field "
	quotedField, ok := strings.CutPrefix(errMsg, prefix)
	if !ok {
		return "", false
	}

	field, err := strconv.Unquote(quotedField)
	if err != nil {
		return "", false
	}
	return field, true
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
