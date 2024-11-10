package main

import (
	"bufio"
	"io"
	"net"
	"os"
)

func (app *application) serve() error {
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
			app.logger.Error("error accepting connection", err.Error())
			os.Exit(1)
		}

		go app.handleConnection(conn)
	}
}

func (app *application) handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	app.logger.Info("accepted connection from", "addr", conn.RemoteAddr())

	for {
		resp, err := Parse(reader)
		if err != nil {
			if err == io.EOF {
				break
			}

			app.logger.Error(err.Error())
			conn.Write([]byte("-ERR internal error\r\n"))
			continue
		}

		app.logger.Debug("Handle command", resp)

		_, err = conn.Write([]byte("+PONG\r\n"))
		if err != nil {
			app.logger.Error(err.Error())
			os.Exit(1)
		}
	}
}
