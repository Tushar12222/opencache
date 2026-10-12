package main

import "net/http"

type Route struct {
	Path    string
	Handler func(w http.ResponseWriter, r *http.Request)
}

var (
	handler = newHandler()
	routes  = [...]Route{
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
