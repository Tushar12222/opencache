package main

import (
	"net/http"
)

type Handler struct {
	Service *Service
}

func newHandler() *Handler {
	service := &Service{}
	return &Handler{
		Service: service,
	}
}

func (hander *Handler) handleGetKeyRequest(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")
	handler.Service.getKey(key)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Received request"))
}

func (hander *Handler) handlePostKeyRequest(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")
	handler.Service.addKey(key)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Received request"))
}

func (hander *Handler) handlePatchKeyRequest(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")
	handler.Service.updateKey(key)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Received request"))
}

func (hander *Handler) handleDeleteKeyRequest(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")
	handler.Service.deleteKey(key)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Received request"))
}
