package main

import (
	"github.com/quic-go/quic-go"
)

type (
	OCPHandlerFunc func(w *quic.Stream, r *OCPRequest)
)

type OCPRoute struct {
	Command string
	Handler OCPHandlerFunc
}
