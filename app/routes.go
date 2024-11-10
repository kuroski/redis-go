package main

import (
	"github.com/codecrafters-io/redis-starter-go/pkg/server"
)

func (app *application) routes() server.CommandHandler {
	handler := NewHandler()
	mux := server.NewServeMux()
	mux.HandleFunc("ping", handler.ping)
	mux.HandleFunc("echo", handler.echo)

	return mux.ServeRESP
}
