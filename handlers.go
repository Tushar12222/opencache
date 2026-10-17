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

func (hander *Handler) handleGetKeyRequest(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")
	value, err := handler.Service.getKey(key)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeResponse(w, http.StatusOK, GetKeyResponse{Message: "key retrieved successfully", Value: value})
}

func (hander *Handler) handlePostKeyRequest(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")
	body := r.Body
	decodedBody, err := JSONDecoder[PostKeyRequest](&body)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	err = handler.Service.addKey(key, decodedBody.Value)
	if err != nil {
		writeErrorResponse(w, http.StatusConflict, err.Error())
		return
	}
	writeResponse(w, http.StatusCreated, PostKeyResponse{Message: "key added successfully"})
}

func (hander *Handler) handlePatchKeyRequest(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")
	body := r.Body
	decodedBody, err := JSONDecoder[PatchKeyRequest](&body)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	err = handler.Service.putKey(key, decodedBody.Value)
	if err != nil {
		writeErrorResponse(w, http.StatusNotFound, err.Error())
		return
	}
	writeResponse(w, http.StatusOK, PatchKeyResponse{Message: "key updated successfully"})
}

func (hander *Handler) handleDeleteKeyRequest(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")
	err := handler.Service.deleteKey(key)
	if err != nil {
		writeErrorResponse(w, http.StatusNotFound, err.Error())
		return
	}
	writeResponse(w, http.StatusOK, DeleteKeyResponse{Message: "key deleted successfully"})
}
