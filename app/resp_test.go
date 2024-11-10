package main

import (
	"bytes"
	"io"
	"log/slog"
	"strings"
	"testing"
)

func (c Command) Equal(other Command) bool {
	if len(c.Args) != len(other.Args) {
		return false
	}

	for i := range c.Args {
		if len(c.Args[i]) != len(other.Args[i]) {
			return false
		}

		for j := range c.Args[i] {
			if c.Args[i][j] != other.Args[i][j] {
				return false
			}
		}
	}

	return bytes.Equal(c.Name, other.Name)
}

func assertEqual(t *testing.T, actual, expected Command) {
	t.Helper()
	if !actual.Equal(expected) {
		var bufferActual bytes.Buffer
		for _, b := range actual.Args {
			bufferActual.Write(b)
		}

		var bufferExpected bytes.Buffer
		for _, b := range expected.Args {
			bufferExpected.Write(b)
		}
		t.Errorf("got: '%v %v'; want: '%v %v'", string(actual.Name), bufferActual.String(), string(expected.Name), bufferExpected.String())
	}
}

func TestCommand(t *testing.T) {
	t.Helper()

	app := &application{
		logger:      slog.New(slog.NewTextHandler(io.Discard, nil)),
		maxBuffSize: 1024,
	}

	tests := []struct {
		input    string
		expected *Command
	}{
		{
			input: "+PING\r\n",
			expected: &Command{
				Name: []byte("PING"),
			},
		},
		{
			input: "*2\r\n$4\r\nECHO\r\n$3\r\nhey\r\n",
			expected: &Command{
				Name: []byte("ECHO"),
				Args: [][]byte{
					[]byte("hey"),
				},
			},
		},
	}

	for index, test := range tests {
		resp := app.NewResp(strings.NewReader(test.input))
		command, err := resp.ReadCommand()
		if err != nil {
			t.Errorf("failed parsing test case %d", index+1)
			t.Fatal(err)
		}

		assertEqual(t, command, *test.expected)
	}
}
