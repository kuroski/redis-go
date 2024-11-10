package main

import (
	"bufio"
	"errors"
	"strconv"
)

type Type byte

const (
	String = '+'
	Bulk   = '$'
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
	case String, Bulk:
	default:
		return nil, errors.New("invalid kind")
	}

	resp.Count, err = strconv.Atoi(string(resp.Data))
	if resp.Type == Bulk {
		if err != nil {
			return nil, errors.New("invalid number of bytes")
		}
	}

	return resp, nil
}

/*
func (r *Resp) Read() (Value, error) {
	_type, err := r.reader.ReadByte()

	if err != nil {
		return Value{}, err
	}

	switch _type {
	case ARRAY:
		return r.readArray()
	case BULK:
		return r.readBulk()
	default:
		fmt.Printf("Unknown type: %v", string(_type))
		return Value{}, nil
	}
}
*/
