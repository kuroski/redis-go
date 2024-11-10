package main

import (
	"log/slog"
	"net"
	"os"
)

type Application struct {
	logger      *slog.Logger
	addr        net.Addr
	maxBuffSize int
	rd          *RespReader
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))

	addr, err := net.ResolveTCPAddr("tcp", net.JoinHostPort("0.0.0.0", "6379"))
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	app := Application{
		logger:      logger,
		addr:        addr,
		maxBuffSize: 1024,
	}

	if err := app.StartServer(); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
