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
	Args [][]byte
}

type RespReader struct {
	rd     *bufio.Reader
	buf    []byte
	logger *slog.Logger
}

func NewReader(rd io.Reader, maxBuffSize int, logger *slog.Logger) *RespReader {
	return &RespReader{
		rd:     bufio.NewReader(rd),
		buf:    make([]byte, maxBuffSize),
		logger: logger,
	}
}

func (rd *RespReader) ReadCommand() (cmd Command, err error) {
	respType, err := rd.rd.ReadByte()
	if err != nil {
		rd.logger.Error(err.Error())
		return Command{}, models.ErrEOF
	}

	switch respType {
	case String, Error: // "+PING\r\n"
		s, _, err := rd.parseString()
		if err != nil {
			return Command{}, err
		}

		cmd.Args = append(cmd.Args, s)
		break
	case Bulk: // "$4\r\nECHO\r\n"
		s, _, err := rd.parseBulk()
		if err != nil {
			return Command{}, err
		}
		cmd.Args = append(cmd.Args, s)
		break
	case Array: // *2\r\n$4\r\nECHO\r\n$3\r\nhey\r\n
		arr, n, err := rd.parseArray()
		if err != nil {
			return Command{}, err
		}

		if n == 0 {
			return Command{}, models.ErrEmptyCommand
		}
		cmd.Args = append(cmd.Args, arr...)
		break
	default:
		return Command{}, models.ErrInvalidDataType.WithArgs(slog.String("respType", string(respType)))
	}

	return cmd, nil
}

func (rd *RespReader) parseLine() (line []byte, n int, err error) {
	line, err = rd.rd.ReadBytes('\n')
	if err != nil {
		return nil, 0, models.ErrMissingCRLFTerminator
	}

	line, found := bytes.CutSuffix(line, []byte{'\r', '\n'})
	if !found {
		return nil, 0, models.ErrMissingCRLFTerminator
	}

	return line, len(line), nil
}

func (rd *RespReader) parseString() (s []byte, n int, err error) {
	line, n, err := rd.parseLine()
	if err != nil {
		return nil, 0, err
	}

	return line, n, nil
}

func (rd *RespReader) parseInteger() (v int, n int, err error) {
	line, n, err := rd.parseLine()
	if err != nil {
		return 0, 0, err
	}

	i64, err := strconv.ParseInt(string(line), 10, 64)
	if err != nil {
		return 0, 0, models.ErrParsingInvalidInteger
	}

	return int(i64), n, nil
}

func (rd *RespReader) parseBulk() (s []byte, n int, err error) {
	size, n, err := rd.parseInteger()
	if err != nil {
		rd.logger.Error(err.Error())
		return nil, 0, models.ErrInvalidBulkFormat
	}

	bulk := make([]byte, size)

	if _, err := rd.rd.Read(bulk); err != nil {
		return nil, 0, err
	}

	rem, n, err := rd.parseLine()
	if err != nil || n != 0 {
		return nil, 0, models.ErrBulkSizeDiffersFromValue.WithArgs(slog.Int("size", size), slog.String("bulk", string(bulk)), slog.String("remaining", string(rem)))
	}

	return bulk, n, nil
}

func (rd *RespReader) parseArray() (arr [][]byte, n int, err error) {
	size, _, err := rd.parseInteger()
	if err != nil {
		rd.logger.Error(err.Error())
		return nil, 0, models.ErrInvalidArrayFormat
	}

	arr = make([][]byte, size)

	for i := 0; i < size; i++ {
		respType, err := rd.rd.ReadByte()
		if err != nil {
			rd.logger.Error(err.Error())
			return nil, 0, models.ErrNotEnoughData
		}

		switch respType {
		case String, Error:
			s, _, err := rd.parseString()
			if err != nil {
				return nil, 0, err
			}
			arr[i] = s
			break
		case Integer:
			i, _, err := rd.parseInteger()
			if err != nil {
				return nil, 0, err
			}
			arr[i] = []byte(strconv.Itoa(i))
			break
		case Bulk:
			s, _, err := rd.parseBulk()
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
