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
	handler.Service.getKey(key)
	writeResponse(w)
}

func (hander *Handler) handlePostKeyRequest(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")
	body := r.Body
	decodedBody, err := JSONDecoder[PostKeyRequest](&body)
	if err != nil {
		writeErrorResponse(w, err)
	}
	handler.Service.addKey(key, decodedBody.Value)
	writeResponse(w)
}

func (hander *Handler) handlePatchKeyRequest(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")
	body := r.Body
	decodedBody, err := JSONDecoder[PatchKeyRequest](&body)
	if err != nil {
		writeErrorResponse(w, err)
		return
	}
	handler.Service.putKey(key, decodedBody.Value)
	writeResponse(w)
}

func (hander *Handler) handleDeleteKeyRequest(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")
	handler.Service.deleteKey(key)
	writeResponse(w)
}
