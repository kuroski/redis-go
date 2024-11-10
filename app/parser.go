package main

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
)

type errProtocol struct {
	msg  string
	args []byte
}

func (err *errProtocol) withArgs(args []byte) *errProtocol {
	err.args = args
	return err
}

func (err *errProtocol) Error() string {
	return fmt.Sprintf("%s '%s'", err.msg, string(err.args))
}

var (
	errNotEnoughData         = errors.New("not enough data")
	errMissingCRLFTerminator = errors.New("missing CRLF terminator")
	errInvalidDataType       = &errProtocol{msg: "invalid data type"}
)

const (
	String  = '+'
	Error   = '-'
	Integer = ':'
	Bulk    = '$'
	Array   = '*'
)

type Parser struct {
}

func Parse(b []byte) (command Command, err error) {
	size := len(b)
	if size <= 1 {
		return Command{}, errNotEnoughData
	}

	if !bytes.HasSuffix(b, []byte("\r\n")) {
		return Command{}, errMissingCRLFTerminator
	}

	respType := b[0]
	data := b[1 : size-2]

	// "*2\r\n$4\r\nECHO\r\n$3\r\nhey\r\n"
	switch respType {
	case String, Error:
		command.Name = data
		break
	case Array:
		var argCount strings.Builder
		for i, p := range data {
			if p == '\r' && data[i+1] == '\n' {
				fmt.Println("FOUND ARRAY COUNT")
				break
			} else {
				argCount.WriteByte(p)
			}
		}
		fmt.Println("======== ARG COUNT ========")
		fmt.Println(argCount.String())
		break
	default:
		return Command{}, errInvalidDataType.withArgs([]byte{respType})
	}

	return command, nil
}
