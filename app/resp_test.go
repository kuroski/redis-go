package main

import (
	"fmt"
	"strconv"
	"testing"
)

func isEmptyRESP(resp Resp) bool {
	return resp.Type == 0 && resp.Count == 0 &&
		resp.Data == nil && resp.Raw == nil
}

func respEquals(a, b Resp) bool {
	if a.Count != b.Count {
		return false
	}
	if a.Type != b.Type {
		return false
	}
	if (a.Data == nil && b.Data != nil) || (a.Data != nil && b.Data == nil) {
		return false
	}
	if string(a.Data) != string(b.Data) {
		return false
	}
	if (a.Raw == nil && b.Raw != nil) || (a.Raw != nil && b.Raw == nil) {
		return false
	}
	if string(a.Raw) != string(b.Raw) {
		return false
	}
	return true
}

func respVOut(a Resp) string {
	var data string
	var raw string
	if a.Data == nil {
		data = "nil"
	} else {
		data = strconv.Quote(string(a.Data))
	}
	if a.Raw == nil {
		raw = "nil"
	} else {
		raw = strconv.Quote(string(a.Raw))
	}
	return fmt.Sprintf("{Type: %d, Count: %d, Data: %s, Raw: %s}",
		a.Type, a.Count, data, raw,
	)
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
		if n != len(test.input) || isEmptyRESP(resp) {
			t.Fatalf("expected good resp")
		}
		if string(resp.Raw) != test.input {
			t.Fatalf("expected '%s', got '%s'", test.input, resp.Raw)
		}
		test.expected.Raw = []byte(test.input)
		switch test.expected.Type {
		case Integer, String, Error:
			test.expected.Data = []byte(test.input[1 : len(test.input)-2])
		}
		if !respEquals(resp, *test.expected) {
			t.Fatalf("expected %v, got %v", respVOut(*test.expected), respVOut(resp))
		}
	}
}
