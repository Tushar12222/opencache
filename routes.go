package main

import "net/http"

type Route struct {
	Path    string
	Handler func(w http.ResponseWriter, r *http.Request)
}

var routes = [...]Route{
	{
		Path:    "GET /value/{key}",
		Handler: handleGetKeyRequest,
	},
}

func RegisterRoutes(server *http.ServeMux) {
	for _, route := range routes {
		server.HandleFunc(route.Path, route.Handler)
	}
}
