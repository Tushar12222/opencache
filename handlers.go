package main

import (
	"net/http"

	"github.com/quic-go/quic-go"
)

type OCPHandler struct {
	Service *Service
}

func newOCPHandler() *OCPHandler {
	cache := NewCache()
	service := NewService(cache)
	return &OCPHandler{
		Service: service,
	}
}

func (handler *OCPHandler) handleOCPKeyExists(w *quic.Stream, r *OCPRequest) {
	key := r.Key
	request := KeyExistsRequest{Key: key}
	if err := request.Validate(); err != nil {
		writeOCPErrorResponse(w, uint8(err.ocpErrorType))
		return
	}
	exists := handler.Service.keyExists(key)

	var valueBytes []byte

	if exists {
		valueBytes = []byte{0x01}
	} else {
		valueBytes = []byte{0x00}
	}

	existsResponseBody, err := toOCPResponseBody(1, 1, valueBytes, nil)
	if err != nil {
		writeOCPErrorResponse(w, uint8(err.ocpErrorType))
		return
	}

	writeOCPResponse(w, 0, existsResponseBody)
}

func (handler *OCPHandler) handleOCPGetKey(w *quic.Stream, r *OCPRequest) {
	key := r.Key
	request := GetKeyRequest{Key: key}
	if err := request.Validate(); err != nil {
		writeOCPErrorResponse(w, uint8(err.ocpErrorType))
		return
	}
	valueType, valueBytes, expiration, err := handler.Service.getKey(key)
	if err != nil {
		writeOCPErrorResponse(w, uint8(err.ocpErrorType))
		return
	}

	if expiration == nil {
		getResponseBody, err := toOCPResponseBody(valueType, uint32(len(valueBytes)), valueBytes, nil)
		if err != nil {
			writeOCPErrorResponse(w, uint8(err.ocpErrorType))
			return
		}
		writeOCPResponse(w, http.StatusOK, getResponseBody)
		return
	}

	getResponseBody, error := toOCPResponseBody(valueType, uint32(len(valueBytes)), valueBytes, expiration)
	if error != nil {
		writeOCPErrorResponse(w, uint8(error.ocpErrorType))
		return
	}
	writeOCPResponse(w, http.StatusOK, getResponseBody)
}

func (handler *OCPHandler) handleOCPPostKey(w *quic.Stream, r *OCPRequest) {
	key := r.Key

	postRequest := PostKeyRequest{
		r.Key,
		r.ValueType,
		r.ValueBytes,
		r.Expiry,
	}

	if err := postRequest.Validate(); err != nil {
		writeOCPErrorResponse(w, uint8(err.ocpErrorType))
		return
	}

	err := handler.Service.addKey(key, postRequest.ValueType, postRequest.ValueBytes, postRequest.Expiry)
	if err != nil {
		writeOCPErrorResponse(w, uint8(err.ocpErrorType))
		return
	}
	writeOCPResponse(w, http.StatusCreated, nil)
}

func (handler *OCPHandler) handleOCPPatchKey(w *quic.Stream, r *OCPRequest) {
	key := r.Key

	patchRequest := PatchKeyRequest{
		Key:        r.Key,
		ValueType:  r.ValueType,
		ValueBytes: r.ValueBytes,
		Expiry:     r.Expiry,
	}

	if err := patchRequest.Validate(); err != nil {
		writeOCPErrorResponse(w, uint8(err.ocpErrorType))
		return
	}

	err := handler.Service.putKey(key, patchRequest.ValueType, patchRequest.ValueBytes, patchRequest.Expiry)
	if err != nil {
		writeOCPErrorResponse(w, uint8(err.ocpErrorType))
		return
	}
	writeOCPResponse(w, http.StatusOK, nil)
}

func (handler *OCPHandler) handleOCPDeleteKey(w *quic.Stream, r *OCPRequest) {
	key := r.Key
	request := DeleteKeyRequest{Key: key}
	if err := request.Validate(); err != nil {
		writeOCPErrorResponse(w, uint8(err.ocpErrorType))
		return
	}
	err := handler.Service.deleteKey(key)
	if err != nil {
		writeOCPErrorResponse(w, uint8(err.ocpErrorType))
		return
	}
	writeOCPResponse(w, http.StatusOK, nil)
}
