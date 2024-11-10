package main

import (
	"bytes"
	"testing"
)

func (c Command) Equal(other Command) bool {
	return bytes.Equal(c.Name, other.Name) && bytes.Equal(c.Args, other.Args)
}

func assertEqual(t *testing.T, actual, expected Command) {
	t.Helper()
	if !actual.Equal(expected) {
		t.Errorf("got: %v; want: %v", actual, expected)
	}
}

func TestParser(t *testing.T) {
	t.Helper()
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
				Args: []byte("hey"),
			},
		},
	}

	for index, test := range tests {
		command, err := Parse([]byte(test.input))
		if err != nil {
			t.Errorf("failed parsing test case %d", index+1)
			t.Fatal(err)
		}

		assertEqual(t, command, *test.expected)
	}
}
