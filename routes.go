package main

import "net/http"

type HandlerFunc func(w http.ResponseWriter, r *http.Request)

type Route struct {
	Path    string
	Handler HandlerFunc
}

var (
	globalHandler = newHandler()
	routes        = [...]Route{
		{
			Path:    "GET /key/{key}/exists",
			Handler: globalHandler.handleKeyExistsRequest,
		},
		{
			Path:    "GET /key/{key}",
			Handler: globalHandler.handleGetKeyRequest,
		},
		{
			Path:    "POST /key/{key}",
			Handler: globalHandler.handlePostKeyRequest,
		},
		{
			Path:    "PATCH /key/{key}",
			Handler: globalHandler.handlePatchKeyRequest,
		},
		{
			Path:    "DELETE /key/{key}",
			Handler: globalHandler.handleDeleteKeyRequest,
		},
	}
)

func RegisterRoutes(server *http.ServeMux) {
	for _, route := range routes {
		server.HandleFunc(route.Path, route.Handler)
	}
}
