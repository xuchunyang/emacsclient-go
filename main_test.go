package main

import (
	"testing"
)

func Test_server_quote_arg(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"'(+ 1 2)", "'(+&_1&_2)"},
		{"(emacs-pid)", "(emacs&-pid)"},
	}
	for _, test := range tests {
		if got := server_quote_arg(test.input); got != test.want {
			t.Errorf("server_quote_arg(%q) = %q, want %q", test.input, got, test.want)
		}
	}
}
