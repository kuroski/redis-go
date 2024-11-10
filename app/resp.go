package main

import (
	"bufio"
	"errors"
)

type Type byte

const (
	String  = '+'
	Error   = '-'
	Integer = ':'
	Bulk    = '$'
	Array   = '*'
)

type Resp struct {
	Type  Type
	Raw   []byte
	Data  []byte
	Count int
}

func Parse(reader *bufio.Reader) (*Resp, error) {
	buff := make([]byte, 1024)
	n, err := reader.Read(buff)
	if err != nil {
		return nil, err
	}

	if n == 0 {
		return nil, errors.New("no data to read")
	}

	if n < 1 {
		return nil, errors.New("not enough data")
	}

	if buff[n-1] != '\n' || buff[n-2] != '\r' {
		return nil, errors.New("invalid resp termination, missing CRLF characters")
	}

	resp := &Resp{
		Type: Type(buff[0]),
		Raw:  buff,
		Data: buff[1 : n-2],
	}

	switch resp.Type {
	//case String, Error, Integer, Bulk, Array:
	case String, Error:
		break
	case Integer:
		if len(resp.Data) == 0 {
			return nil, errors.New("invalid integer")
		}

		var i int
		if resp.Data[0] == '-' {
			if len(resp.Data) == 1 {
				return nil, errors.New("invalid negative integer")
			}
			i++
		}

		for ; i < len(resp.Data); i++ {
			if resp.Data[i] < '0' || resp.Data[i] > '9' {
				return nil, errors.New("invalid integer")
			}
		}
		break
	default:
		return nil, errors.New("invalid kind")
	}

	return resp, nil
}
