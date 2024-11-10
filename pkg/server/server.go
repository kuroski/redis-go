package server

import (
	"fmt"
	"github.com/codecrafters-io/redis-starter-go/pkg/resp"
	"log"
	"net"
	"strings"
)

type CommandHandler func(conn net.Conn, cmd resp.Command)

type ServeMux struct {
	handlers map[string]CommandHandler
}

func NewServeMux() *ServeMux {
	return &ServeMux{
		handlers: make(map[string]CommandHandler),
	}
}

func (mux *ServeMux) HandleFunc(command string, handler CommandHandler) {
	if handler == nil {
		panic("redcon: nil handler")
	}
	mux.handlers[strings.ToLower(command)] = handler
}

func (mux *ServeMux) ServeRESP(conn net.Conn, cmd resp.Command) {
	if handler, ok := mux.handlers[strings.ToLower(string(cmd.Args[0]))]; ok {
		handler(conn, cmd)
	} else {
		log.Printf("command not supported '%s'", cmd.Args)
		_, _ = conn.Write([]byte(fmt.Sprintf("-ERR unknown command '%s'\r\n", cmd.Args[0])))
	}
}
