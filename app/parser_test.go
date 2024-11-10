package main

import (
	"testing"
)

type Type byte

const (
	String  Type = '+'
	Error   Type = '-'
	Integer Type = ':'
	Bulk    Type = '$'
	Array   Type = '*'
)

type Resp struct {
	Type Type
	Data []byte
}

func ReadNextResp(b []byte) (n int, resp Resp) {
	return 0, Resp{}
}

func TestResp(t *testing.T) {
	t.Helper()
	tests := []struct {
		input    string
		expected *Resp
		err      bool
	}{
		{
			input: "+OK\r\n",
			expected: &Resp{
				Type: String,
				Data: []byte("OK"),
			},
			err: false,
		},
	}

	for _, test := range tests {
		n, resp := ReadNextResp([]byte(test.input))
		t.Log("-----------------")
		t.Log(n, resp)
	}
}
