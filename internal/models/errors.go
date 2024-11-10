package models

import (
	"errors"
	"fmt"
)

var (
	ErrNotEnoughData         = errors.New("not enough data")
	ErrMissingCRLFTerminator = errors.New("missing CRLF terminator")
	ErrInvalidStringFormat   = errors.New("invalid string format")
	ErrInvalidDataType       = &ErrProtocol{msg: "invalid data type"}
)

type ErrProtocol struct {
	msg  string
	args []byte
}

func (err *ErrProtocol) WithArgs(args []byte) *ErrProtocol {
	err.args = args
	return err
}

func (err *ErrProtocol) Error() string {
	return fmt.Sprintf("%s '%s'", err.msg, string(err.args))
}
