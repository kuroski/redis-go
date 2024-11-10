package main

import (
	"bufio"
	"bytes"
	"github.com/codecrafters-io/redis-starter-go/internal/models"
	"io"
	"log/slog"
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

	// "*2\r\n$4\r\nECHO\r\n$3\r\nhey\r\n"
	// "+PING\r\n"
	switch respType {
	case String, Error:
		s, _, err := r.readString()
		if err != nil {
			return Command{}, models.ErrInvalidStringFormat
		}
		command.Name = s
		break
	case Array:
		//var argCount strings.Builder
		//for i, p := range data {
		//	if p == '\r' && data[i+1] == '\n' {
		//		fmt.Println("FOUND ARRAY COUNT")
		//		break
		//	} else {
		//		argCount.WriteByte(p)
		//	}
		//}
		//fmt.Println("======== ARG COUNT ========")
		//fmt.Println(argCount.String())
		break
	default:
		return Command{}, models.ErrInvalidDataType.WithArgs([]byte{respType})
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

	return line, n, err
}

//func Parse(b []byte) (command Command, err error) {
//	size := len(b)
//	if size <= 1 {
//		return Command{}, models.ErrNotEnoughData
//	}
//
//	if !bytes.HasSuffix(b, []byte("\r\n")) {
//		return Command{}, models.ErrMissingCRLFTerminator
//	}
//
//	respType := b[0]
//	data := b[1 : size-2]
//
//	// "*2\r\n$4\r\nECHO\r\n$3\r\nhey\r\n"
//	switch respType {
//	case String, Error:
//		command.Name = data
//		break
//	case Array:
//		var argCount strings.Builder
//		for i, p := range data {
//			if p == '\r' && data[i+1] == '\n' {
//				fmt.Println("FOUND ARRAY COUNT")
//				break
//			} else {
//				argCount.WriteByte(p)
//			}
//		}
//		fmt.Println("======== ARG COUNT ========")
//		fmt.Println(argCount.String())
//		break
//	default:
//		return Command{}, models.ErrInvalidDataType.WithArgs([]byte{respType})
//	}
//
//	return command, nil
//}

//func parseArray(b []byte) {
//	var argCount strings.Builder
//	for i, p := range b {
//		if p == '\r' && b[i+1] == '\n' {
//			fmt.Println("FOUND ARRAY COUNT")
//			break
//		} else {
//			argCount.WriteByte(p)
//		}
//	}
//}

//func parseInteger() (x int, n int, err error) {
//	line, n, err := r.readLine()
//	if err != nil {
//		return 0, 0, err
//	}
//	i64, err := strconv.ParseInt(string(line), 10, 64)
//	if err != nil {
//		return 0, n, err
//	}
//	return int(i64), n, nil
//}
