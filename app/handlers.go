package main

import (
	"fmt"
	"github.com/codecrafters-io/redis-starter-go/pkg/resp"
	"net"
	"sync"
)

type Handler struct {
	itemsMux sync.RWMutex
	items    map[string][]byte
}

func NewHandler() *Handler {
	return &Handler{
		items: make(map[string][]byte),
	}
}

func (h *Handler) ping(conn net.Conn, cmd resp.Command) {
	_, err := conn.Write([]byte("+PONG\r\n"))
	if err != nil {
		conn.Write([]byte(fmt.Sprintf("-ERR %s\r\n", err.Error())))
	}
}

func (h *Handler) echo(conn net.Conn, cmd resp.Command) {
	_, err := conn.Write([]byte(fmt.Sprintf("+%s\r\n", cmd.Args[1])))
	if err != nil {
		conn.Write([]byte(fmt.Sprintf("-ERR %s\r\n", err.Error())))
	}
}
