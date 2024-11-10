package main

import (
	"bytes"
	"errors"
	"fmt"
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

	switch respType {
	case String, Error:
		command.Name = data
		break
	default:
		return Command{}, errInvalidDataType.withArgs([]byte{respType})
	}

	return command, nil
}
