package models

import (
	"errors"
	"fmt"
	"log/slog"
)

var (
	ErrNotEnoughData            = errors.New("not enough data")
	ErrMissingCRLFTerminator    = errors.New("missing CRLF terminator")
	ErrParsingInvalidInteger    = errors.New("trying to parse invalid integer")
	ErrInvalidBulkFormat        = errors.New("invalid bulk format")
	ErrBulkSizeDiffersFromValue = &ErrProtocol{msg: "bulk size is different from the actual value"}
	ErrInvalidDataType          = &ErrProtocol{msg: "invalid data type"}
)

type ErrProtocol struct {
	msg  string
	args []any
}

func (err *ErrProtocol) WithArgs(args ...any) *ErrProtocol {
	err.args = args
	return err
}

func (err *ErrProtocol) Error() string {
	args := slog.Group("args", err.args...)
	return fmt.Sprint(err.msg, " ", args)
	//return fmt.Sprintf("%s '%s'", err.msg, args.String())
}
