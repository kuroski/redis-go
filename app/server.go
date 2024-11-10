package main

import (
	"errors"
	"fmt"
	"github.com/codecrafters-io/redis-starter-go/internal/models"
	"net"
	"os"
	"strings"
)

func (app *Application) StartServer() error {
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

		go app.handleConnection(conn)
	}
}

func (app *Application) handleConnection(conn net.Conn) {
	defer conn.Close()
	app.logger.Info("accepted connection from", "addr", conn.RemoteAddr())

	for {
		rd := NewReader(conn, app.maxBuffSize, app.logger)
		cmd, err := rd.ReadCommand()
		if err != nil {
			if errors.Is(err, models.ErrEOF) {
				break
			}

			app.logger.Error(err.Error())
			conn.Write([]byte(fmt.Sprintf("-ERR %s\r\n", err.Error())))
			continue
		}

		switch strings.ToLower(string(cmd.Args[0])) {
		case "ping":
			_, err = conn.Write([]byte("+PONG\r\n"))
			if err != nil {
				app.logger.Error(err.Error())
			}
		case "echo":
			_, err = conn.Write([]byte(fmt.Sprintf("+%s\r\n", cmd.Args[1])))
			if err != nil {
				app.logger.Error(err.Error())
			}
		default:
			app.logger.Info(fmt.Sprintf("command not supported '%s'", cmd.Args))
		}
	}
}
