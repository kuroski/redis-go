package main

import "bytes"

type Type byte

const (
	String  Type = '+'
	Error   Type = '-'
	Integer Type = ':'
	Bulk    Type = '$'
	Array   Type = '*'
)

type Resp struct {
	Type  Type   // resp data type
	Data  []byte // the actual data of the resp, "$5\r\nhello\r\n" == byte["hello"]
	Raw   []byte // the raw value of the payload
	Count int    // length of bulk strings/errors
}

func ReadNextResp(b []byte) (n int, resp Resp) {
	size := len(b)
	if size <= 1 || !bytes.HasSuffix(b, []byte("\r\n")) {
		return 0, Resp{} // not enough data || missing CRLF terminator
	}

	resp.Type = Type(b[0])
	resp.Raw = bytes.Clone(b)
	resp.Data = b[1 : size-2]

	switch resp.Type {
	case String, Error:
		break
	default:
		return 0, Resp{} // invalid data type
	}

	return len(resp.Raw), resp
}
