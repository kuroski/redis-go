package main

import (
	"github.com/codecrafters-io/redis-starter-go/pkg/server"
)

func (app *application) routes() server.Handler {
	handler := NewHandler()
	mux := server.NewServeMux()
	mux.HandleFunc("ping", handler.ping)
	mux.HandleFunc("echo", handler.echo)
	mux.HandleFunc("get", handler.get)
	mux.HandleFunc("set", handler.set)

	return mux.ServeRESP
}
