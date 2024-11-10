package main

import (
	"fmt"
	"io"
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
	buf := make([]byte, 1024)

	app.logger.Info("accepted connection from", "addr", conn.RemoteAddr())

	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}

			app.logger.Error(err.Error())
			conn.Write([]byte(fmt.Sprintf("-ERR %s\r\n", err.Error())))
			continue
		}

		command := buf[:n]
		n, resp := ReadNextResp(command)
		//if err != nil {
		//	app.logger.Error(err.Error())
		//	conn.Write([]byte(fmt.Sprintf("-ERR %s\r\n", err.Error())))
		//	continue
		//}
		//
		//// TODO: Handle Writer later
		//app.logger.Debug("Handle command", resp)
		//
		//_, err = conn.Write([]byte("+PONG\r\n"))
		//if err != nil {
		//	app.logger.Error(err.Error())
		//	os.Exit(1)
		//}
	}
}

func readCommands() {

}
