package main

import (
	"log/slog"
	"net"
	"os"
)

type application struct {
	logger *slog.Logger
	addr   net.Addr
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))

	addr, err := net.ResolveTCPAddr("tcp", net.JoinHostPort("0.0.0.0", "6379"))
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	app := application{
		logger: logger,
		addr:   addr,
	}

	if err := app.Serve(); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
