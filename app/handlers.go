package main

import (
	"fmt"
	"github.com/codecrafters-io/redis-starter-go/pkg/resp"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Item struct {
	Value  []byte
	Expiry time.Time
}

type Handler struct {
	itemsMux sync.RWMutex
	items    map[string]*Item
}

func NewHandler() *Handler {
	return &Handler{
		items: make(map[string]*Item),
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
	defer h.itemsMux.Unlock()

	item, ok := h.items[string(cmd.Args[1])]

	if !ok {
		conn.Write([]byte(fmt.Sprint("$-1\r\n"))) // null bulk string
		return
	}

	now := time.Now()
	if !item.Expiry.IsZero() && item.Expiry.Before(now) {
		conn.Write([]byte(fmt.Sprint("$-1\r\n"))) // null bulk string
		return
	}

	n := len(item.Value)
	_, err := conn.Write([]byte(fmt.Sprintf("$%d\r\n%s\r\n", n, item.Value))) // null bulk string
	if err != nil {
		conn.Write([]byte(fmt.Sprintf("-ERR %s\r\n", err.Error())))
	}
}

func (h *Handler) set(conn net.Conn, cmd resp.Command) {
	if len(cmd.Args) < 3 {
		conn.Write([]byte(fmt.Sprint("-ERR wrong number of arguments for 'set' command\r\n")))
		return
	}

	key := string(cmd.Args[1])
	value := cmd.Args[2]
	var exp time.Time

	if len(cmd.Args) > 3 {
		switch strings.ToLower(string(cmd.Args[3])) {
		case "px":
			expMillis, err := strconv.Atoi(string(cmd.Args[4]))
			if err != nil {
				conn.Write([]byte(fmt.Sprint("-ERR wrong number of arguments for 'set' command\r\n")))
			}
			exp = time.Now().Add(time.Millisecond * time.Duration(expMillis))
		default:
			conn.Write([]byte(fmt.Sprint("-ERR wrong number of arguments for 'set' command\r\n")))
		}
	}

	h.itemsMux.Lock()
	h.items[key] = &Item{
		Value:  value,
		Expiry: exp,
	}
	h.itemsMux.Unlock()

	_, err := conn.Write([]byte("+OK\r\n"))
	if err != nil {
		conn.Write([]byte(fmt.Sprintf("-ERR %s\r\n", err.Error())))
	}
}
