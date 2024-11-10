package main

import (
	"bufio"
)

type Type byte

const (
	STRING = '+'
	BULK   = '$'
)

type Resp struct {
	Type  Type
	Raw   []byte
	Data  []byte
	Count int
}

func Parse(reader *bufio.Reader) (*Resp, error) {
	return &Resp{}, nil
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
