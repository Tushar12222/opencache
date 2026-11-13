package main

import (
	"net/http"

	"github.com/quic-go/quic-go"
)

type Handler struct {
	Service *Service
}

type OCPHandler struct {
	Service *Service
}

func newHandler() *Handler {
	cache := NewCache()
	service := NewService(cache)
	return &Handler{
		Service: service,
	}
}

func newOCPHandler() *OCPHandler {
	cache := NewCache()
	service := NewService(cache)
	return &OCPHandler{
		Service: service,
	}
}

func GetGlobalHandler() *Handler {
	return globalHandler
}

func (handler *OCPHandler) handleOCPKeyExists(w *quic.Stream, r *OCPRequest) {
	key := r.Key
	request := KeyExistsRequest{Key: key}
	if err := request.Validate(); err != nil {
		writeOCPErrorResponse(w, uint8(BadRequestOCPError))
		return
	}
	exists := handler.Service.keyExists(key)
	writeOCPResponse(w, http.StatusOK, &GetKeyResponse{Message: "key retrieved successfully", Data: map[string]any{
		"value": exists,
	}})
}

func (handler *OCPHandler) handleOCPGetKey(w *quic.Stream, r *OCPRequest) {
	key := r.Key
	request := GetKeyRequest{Key: key}
	if err := request.Validate(); err != nil {
		writeOCPErrorResponse(w, uint8(BadRequestOCPError))
		return
	}
	val, expiration, err := handler.Service.getKey(key)
	if err != nil {
		writeOCPErrorResponse(w, uint8(KeyNotFoundOCPError))
		return
	}

	if expiration == nil {
		writeOCPResponse(w, http.StatusOK, &GetKeyResponse{Message: "key retrieved successfully", Data: map[string]any{
			"value": val,
		}})
		return
	}

	writeOCPResponse(w, http.StatusOK, &GetKeyResponse{Message: "key retrieved successfully", Data: map[string]any{
		"value": val,
		"exp":   *expiration,
	}})
}

func (handler *OCPHandler) handleOCPPostKey(w *quic.Stream, r *OCPRequest) {
	key := r.Key

	postRequest := PostKeyRequest{
		r.Key,
		r.Value,
		0,
	}

	if err := postRequest.Validate(); err != nil {
		writeOCPErrorResponse(w, uint8(BadRequestOCPError))
		return
	}

	err := handler.Service.addKey(key, postRequest.Value, postRequest.Expiry)
	if err != nil {
		writeOCPErrorResponse(w, uint8(DuplicateKeyOCPError))
		return
	}
	writeOCPResponse(w, http.StatusCreated, &PostKeyResponse{Message: "key added successfully"})
}

func (handler *OCPHandler) handleOCPPatchKey(w *quic.Stream, r *OCPRequest) {
	key := r.Key

	patchRequest := PatchKeyRequest{
		Key:   r.Key,
		Value: r.Value,
	}

	if err := patchRequest.Validate(); err != nil {
		writeOCPErrorResponse(w, uint8(BadRequestOCPError))
		return
	}

	err := handler.Service.putKey(key, patchRequest.Value, patchRequest.Expiry)
	if err != nil {
		writeOCPErrorResponse(w, uint8(KeyNotFoundOCPError))
		return
	}
	writeOCPResponse(w, http.StatusOK, &PatchKeyResponse{Message: "key updated successfully"})
}

func (handler *OCPHandler) handleOCPDeleteKey(w *quic.Stream, r *OCPRequest) {
	key := r.Key
	request := DeleteKeyRequest{Key: key}
	if err := request.Validate(); err != nil {
		writeOCPErrorResponse(w, uint8(BadRequestOCPError))
		return
	}
	err := handler.Service.deleteKey(key)
	if err != nil {
		writeOCPErrorResponse(w, uint8(KeyNotFoundOCPError))
		return
	}
	writeOCPResponse(w, http.StatusOK, &DeleteKeyResponse{Message: "key deleted successfully"})
}

func (handler *Handler) handleKeyExistsRequest(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")
	request := KeyExistsRequest{Key: key}
	if err := request.Validate(); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	exists := handler.Service.keyExists(key)
	writeResponse(w, http.StatusOK, GetKeyResponse{Message: "key retrieved successfully", Data: map[string]any{
		"value": exists,
	}})
}

func (handler *Handler) handleGetKeyRequest(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")
	request := GetKeyRequest{Key: key}
	if err := request.Validate(); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	val, expiration, err := handler.Service.getKey(key)
	if err != nil {
		writeErrorResponse(w, http.StatusNotFound, err.Error())
		return
	}

	if expiration == nil {
		writeResponse(w, http.StatusOK, GetKeyResponse{Message: "key retrieved successfully", Data: map[string]any{
			"value": val,
		}})
		return
	}

	writeResponse(w, http.StatusOK, GetKeyResponse{Message: "key retrieved successfully", Data: map[string]any{
		"value": val,
		"exp":   *expiration,
	}})
}

func (handler *Handler) handlePostKeyRequest(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")
	body := r.Body

	decodedRequest, err := JSONDecoder[PostKeyRequest](&body)
	if err != nil {
		parsedErr := parseDecodeError(err)
		writeErrorResponse(w, http.StatusBadRequest, parsedErr.Error())
		return
	}

	decodedRequest.Key = key
	if err = decodedRequest.Validate(); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	err = handler.Service.addKey(key, decodedRequest.Value, decodedRequest.Expiry)
	if err != nil {
		writeErrorResponse(w, http.StatusConflict, err.Error())
		return
	}
	writeResponse(w, http.StatusCreated, PostKeyResponse{Message: "key added successfully"})
}

func (handler *Handler) handlePatchKeyRequest(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")
	body := r.Body
	decodedRequest, err := JSONDecoder[PatchKeyRequest](&body)
	if err != nil {
		parsedErr := parseDecodeError(err)
		writeErrorResponse(w, http.StatusBadRequest, parsedErr.Error())
		return
	}

	decodedRequest.Key = key
	if err = decodedRequest.Validate(); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	err = handler.Service.putKey(key, decodedRequest.Value, decodedRequest.Expiry)
	if err != nil {
		writeErrorResponse(w, http.StatusNotFound, err.Error())
		return
	}
	writeResponse(w, http.StatusOK, PatchKeyResponse{Message: "key updated successfully"})
}

func (handler *Handler) handleDeleteKeyRequest(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")
	request := DeleteKeyRequest{Key: key}
	if err := request.Validate(); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	err := handler.Service.deleteKey(key)
	if err != nil {
		writeErrorResponse(w, http.StatusNotFound, err.Error())
		return
	}
	writeResponse(w, http.StatusOK, DeleteKeyResponse{Message: "key deleted successfully"})
}
