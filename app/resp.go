package main

import (
	"bufio"
	"bytes"
	"github.com/codecrafters-io/redis-starter-go/internal/models"
	"io"
	"log/slog"
	"strconv"
)

type Command struct {
	Name []byte
	Args []byte
}

type Resp struct {
	logger *slog.Logger
	reader *bufio.Reader
}

func (app *application) NewResp(rd io.Reader) *Resp {
	return &Resp{
		reader: bufio.NewReaderSize(rd, app.maxBuffSize),
		logger: app.logger,
	}
}

const (
	String  = '+'
	Error   = '-'
	Integer = ':'
	Bulk    = '$'
	Array   = '*'
)

func (r *Resp) ReadCommand() (command Command, err error) {
	respType, err := r.reader.ReadByte()
	if err != nil {
		r.logger.Error(err.Error())
		return Command{}, models.ErrNotEnoughData
	}

	switch respType {
	case String, Error: // "+PING\r\n"
		s, _, err := r.readString()
		if err != nil {
			return Command{}, err
		}
		command.Name = s
		break
	case Bulk: // "$4\r\nECHO\r\n"
		s, _, err := r.readBulk()
		if err != nil {
			return Command{}, err
		}
		command.Name = s
		break
	//case Array:
	//	arr, _, err := r.readArray()
	//	if err != nil {
	//		return Command{}, models.ErrInvalidStringFormat
	//	}
	//	command.Name = s
	//	break
	//	//var argCount strings.Builder
	//	//for i, p := range data {
	//	//	if p == '\r' && data[i+1] == '\n' {
	//	//		fmt.Println("FOUND ARRAY COUNT")
	//	//		break
	//	//	} else {
	//	//		argCount.WriteByte(p)
	//	//	}
	//	//}
	//	//fmt.Println("======== ARG COUNT ========")
	//	//fmt.Println(argCount.String())
	//	break
	default:
		return Command{}, models.ErrInvalidDataType.WithArgs(slog.String("respType", string(respType)))
	}

	return command, nil
}

func (r *Resp) readLine() (line []byte, n int, err error) {
	line, err = r.reader.ReadBytes('\n')
	if err != nil {
		return nil, 0, models.ErrMissingCRLFTerminator
	}

	line, found := bytes.CutSuffix(line, []byte{'\r', '\n'})
	if !found {
		return nil, 0, models.ErrMissingCRLFTerminator
	}

	return line, len(line), nil
}

func (r *Resp) readString() (s []byte, n int, err error) {
	line, n, err := r.readLine()
	if err != nil {
		return nil, 0, err
	}

	return line, n, nil
}

func (r *Resp) readInteger() (v int, n int, err error) {
	line, n, err := r.readLine()
	if err != nil {
		return 0, 0, err
	}

	i64, err := strconv.ParseInt(string(line), 10, 64)
	if err != nil {
		return 0, 0, models.ErrParsingInvalidInteger
	}

	return int(i64), n, nil
}

func (r *Resp) readBulk() (s []byte, n int, err error) {
	size, n, err := r.readInteger()
	if err != nil {
		r.logger.Error(err.Error())
		return nil, 0, models.ErrInvalidBulkFormat
	}

	bulk := make([]byte, size)

	if _, err := r.reader.Read(bulk); err != nil {
		return nil, 0, err
	}

	rem, n, err := r.readLine()
	if err != nil || n != 0 {
		return nil, 0, models.ErrBulkSizeDiffersFromValue.WithArgs(slog.Int("size", size), slog.String("bulk", string(bulk)), slog.String("remaining", string(rem)))
	}

	return bulk, n, nil
}
