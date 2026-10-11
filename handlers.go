package main

import (
	"fmt"
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
	fmt.Println("Received key: ", key)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Received request"))
}

func (hander *Handler) handlePostKeyRequest(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")
	fmt.Println("Received key: ", key)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Received request"))
}

func (hander *Handler) handlePatchKeyRequest(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")
	fmt.Println("Received key: ", key)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Received request"))
}

func (hander *Handler) handleDeleteKeyRequest(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")
	fmt.Println("Received key: ", key)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Received request"))
}
