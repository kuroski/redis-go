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

func (h *Handler) get(conn net.Conn, cmd resp.Command) {
	if len(cmd.Args) != 2 {
		conn.Write([]byte(fmt.Sprint("-ERR wrong number of arguments for 'get' command\r\n")))
		return
	}

	h.itemsMux.Lock()
	value, ok := h.items[string(cmd.Args[1])]
	h.itemsMux.Unlock()

	if !ok {
		conn.Write([]byte(fmt.Sprint("$-1\r\n"))) // null bulk string
		return
	}

	n := len(value)
	_, err := conn.Write([]byte(fmt.Sprintf("$%d\r\n%s\r\n", n, value))) // null bulk string
	if err != nil {
		conn.Write([]byte(fmt.Sprintf("-ERR %s\r\n", err.Error())))
	}
}

func (h *Handler) set(conn net.Conn, cmd resp.Command) {
	if len(cmd.Args) != 3 {
		conn.Write([]byte(fmt.Sprint("-ERR wrong number of arguments for 'set' command\r\n")))
		return
	}

	h.itemsMux.Lock()
	h.items[string(cmd.Args[1])] = cmd.Args[2]
	h.itemsMux.Unlock()

	_, err := conn.Write([]byte("+OK\r\n"))
	if err != nil {
		conn.Write([]byte(fmt.Sprintf("-ERR %s\r\n", err.Error())))
	}
}
