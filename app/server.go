package main

import (
	"log/slog"
	"net"
	"os"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))

	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		logger.Error("Failed to bind to port 6379")
		os.Exit(1)
	}

	conn, err := l.Accept()
	if err != nil {
		logger.Error("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
	defer conn.Close()

	buff := make([]byte, 128)

	_, err = conn.Read(buff)
	if err != nil {
		logger.Error("Error reading command: ", err.Error())
		os.Exit(1)
	}

	_, err = conn.Write([]byte("+PONG\r\n"))
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
