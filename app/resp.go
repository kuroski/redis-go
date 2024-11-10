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
	Name string
	Args [][]byte
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

// Reader

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
		command.Name = string(s)
		break
	case Bulk: // "$4\r\nECHO\r\n"
		s, _, err := r.readBulk()
		if err != nil {
			return Command{}, err
		}
		command.Name = string(s)
		break
	case Array: // *2\r\n$4\r\nECHO\r\n$3\r\nhey\r\n
		arr, n, err := r.readArray()
		if err != nil {
			return Command{}, err
		}

		if n == 0 {
			return Command{}, models.ErrEmptyCommand
		}
		command.Name = string(arr[0])
		command.Args = arr[1:n]
		break
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

func (r *Resp) readArray() (arr [][]byte, n int, err error) {
	size, _, err := r.readInteger()
	if err != nil {
		r.logger.Error(err.Error())
		return nil, 0, models.ErrInvalidArrayFormat
	}

	arr = make([][]byte, size)

	for i := 0; i < size; i++ {
		respType, err := r.reader.ReadByte()
		if err != nil {
			r.logger.Error(err.Error())
			return nil, 0, models.ErrNotEnoughData
		}

		switch respType {
		case String, Error:
			s, _, err := r.readString()
			if err != nil {
				return nil, 0, err
			}
			arr[i] = s
			break
		case Integer:
			i, _, err := r.readInteger()
			if err != nil {
				return nil, 0, err
			}
			arr[i] = []byte(strconv.Itoa(i))
			break
		case Bulk:
			s, _, err := r.readBulk()
			if err != nil {
				return nil, 0, err
			}
			arr[i] = s
			break
		default:
			return nil, 0, models.ErrInvalidDataType.WithArgs(slog.String("respType", string(respType)))
		}
	}

	return arr, size, nil
}
