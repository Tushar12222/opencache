package main

import (
	"net/http"
)

type Handler struct {
	Service *Service
}

func newHandler() *Handler {
	cache := NewCache()
	service := NewService(cache)
	return &Handler{
		Service: service,
	}
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
	value, err := handler.Service.getKey(key)
	if err != nil {
		writeErrorResponse(w, http.StatusNotFound, err.Error())
		return
	}
	writeResponse(w, http.StatusOK, GetKeyResponse{Message: "key retrieved successfully", Data: map[string]any{
		"value": value,
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
