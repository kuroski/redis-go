package main

import (
	"errors"
	"fmt"
	"github.com/codecrafters-io/redis-starter-go/pkg/resp"
	"github.com/codecrafters-io/redis-starter-go/pkg/server"
	"io"
	"log/slog"
	"net"
	"os"
)

type application struct {
	logger      *slog.Logger
	addr        net.Addr
	maxBuffSize int
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))

	addr, err := net.ResolveTCPAddr("tcp", net.JoinHostPort("0.0.0.0", "6379"))
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	app := application{
		logger:      logger,
		addr:        addr,
		maxBuffSize: 1024,
	}

	if err := app.startServer(app.routes()); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}

func (app *application) startServer(handler server.Handler) error {
	l, err := net.Listen(app.addr.Network(), app.addr.String())
	if err != nil {
		app.logger.Error("failed to bind to", "addr", app.addr)
		os.Exit(1)
	}

	defer l.Close()
	app.logger.Info("serving", "addr", app.addr.String())

	for {
		conn, err := l.Accept()
		if err != nil {
			app.logger.Error("error accepting connection", "err", err.Error())
			os.Exit(1)
		}

		go app.handleConnection(conn, handler)
	}
}

func (app *application) handleConnection(conn net.Conn, handler server.Handler) {
	defer conn.Close()
	app.logger.Info("accepted connection from", "addr", conn.RemoteAddr())

	for {
		rd := resp.NewReader(conn, app.maxBuffSize)
		cmd, err := rd.ReadCommand()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break // close connection, otherwise it will enter on an infinite loop
			}

			app.logger.Error(err.Error())
			conn.Write([]byte(fmt.Sprintf("-ERR %s\r\n", err.Error())))
			continue
		}

		handler(conn, cmd)
	}
}
