package main

import "net/http"

type HandlerFunc func(w http.ResponseWriter, r *http.Request)

type Route struct {
	Path    string
	Handler HandlerFunc
}

var (
	handler = newHandler()
	routes  = [...]Route{
		{
			Path:    "GET /key/{key}/exists",
			Handler: handler.handleKeyExistsRequest,
		},
		{
			Path:    "GET /key/{key}",
			Handler: handler.handleGetKeyRequest,
		},
		{
			Path:    "POST /key/{key}",
			Handler: handler.handlePostKeyRequest,
		},
		{
			Path:    "PATCH /key/{key}",
			Handler: handler.handlePatchKeyRequest,
		},
		{
			Path:    "DELETE /key/{key}",
			Handler: handler.handleDeleteKeyRequest,
		},
	}
)

func RegisterRoutes(server *http.ServeMux) {
	for _, route := range routes {
		server.HandleFunc(route.Path, route.Handler)
	}
}
