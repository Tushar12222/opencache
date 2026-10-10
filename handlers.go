package main

import (
	"fmt"
	"net/http"
)

func handleGetKeyRequest(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")
	fmt.Println("Received key: ", key)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Received request"))
}
