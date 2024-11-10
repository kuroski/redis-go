package main

import (
	"fmt"
	"net"
	"os"
)

func (app *application) Serve() error {
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

func (app *application) handleConnection(conn net.Conn) {
	defer conn.Close()

	app.logger.Info("accepted connection from", "addr", conn.RemoteAddr())

	for {
		resp := app.NewResp(conn)
		cmd, err := resp.ReadCommand()
		if err != nil {
			app.logger.Error(err.Error())
			conn.Write([]byte(fmt.Sprintf("-ERR %s\r\n", err.Error())))
			continue
		}

		app.logger.Debug("Handle command", resp)

		switch cmd.Name {
		case "PING":
			_, err = conn.Write([]byte("+PONG\r\n"))
			if err != nil {
				app.logger.Error(err.Error())
				os.Exit(1)
			}
		case "ECHO":
			_, err = conn.Write([]byte(fmt.Sprintf("+%s\r\n", cmd.Args)))
			if err != nil {
				app.logger.Error(err.Error())
				os.Exit(1)
			}
		default:
			app.logger.Info(fmt.Sprintf("command not supported '%s'", cmd.Name))
			os.Exit(1)
		}
	}
}
