package main

import (
	"io"
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

	for {
		conn, err := l.Accept()
		if err != nil {
			logger.Error("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go func() {
			for {
				buff := make([]byte, 1024)

				_, err = conn.Read(buff)
				if err != nil {
					if err == io.EOF {
						break
					}
					logger.Error("Error reading from the client ", err.Error())
					os.Exit(1)
				}
				_, err = conn.Write([]byte("+PONG\r\n"))
				if err != nil {
					logger.Error(err.Error())
					os.Exit(1)
				}
			}

		}()
	}
}
