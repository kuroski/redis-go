package main

type Type byte

type Resp struct {
	Type  Type
	Raw   []byte
	Data  []byte
	Count int
}

const (
	String  = '+'
	Error   = '-'
	Integer = ':'
	Bulk    = '$'
	Array   = '*'
)
