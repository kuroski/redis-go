package main

import (
	"bytes"
	"io"
	"log/slog"
	"strings"
	"testing"
)

func (cmd *Command) Equal(other Command) bool {
	if len(cmd.Args) != len(other.Args) {
		return false
	}

	for i := range cmd.Args {
		if len(cmd.Args[i]) != len(other.Args[i]) {
			return false
		}

		for j := range cmd.Args[i] {
			if cmd.Args[i][j] != other.Args[i][j] {
				return false
			}
		}
	}

	return true
}

func (cmd *Command) String() string {
	var bufferActual bytes.Buffer
	for _, b := range cmd.Args {
		bufferActual.Write(b)
	}
	return bufferActual.String()
}

func assertEqual(t *testing.T, actual, expected Command) {
	t.Helper()
	if !actual.Equal(expected) {
		t.Errorf("got: '%v'; want: '%v'", actual, expected)
	}
}

func TestReadCommand(t *testing.T) {
	t.Helper()

	app := &Application{
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
				Args: [][]byte{
					[]byte("PING"),
				},
			},
		},
		{
			input: "*2\r\n$4\r\nECHO\r\n$3\r\nhey\r\n",
			expected: &Command{
				Args: [][]byte{
					[]byte("ECHO"),
					[]byte("hey"),
				},
			},
		},
	}

	for index, test := range tests {
		rd := app.NewReader(strings.NewReader(test.input))
		command, err := rd.ReadCommand()
		if err != nil {
			t.Errorf("failed parsing test case %d", index+1)
			t.Fatal(err)
		}

		assertEqual(t, command, *test.expected)
	}
}
