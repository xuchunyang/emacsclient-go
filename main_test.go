package main

import (
	"testing"
)

func TestQuoteAndUnquote(t *testing.T) {
	tests := []struct {
		unquoted string
		quoted  string
	}{
		{"'(+ 1 2)", "'(+&_1&_2)"},
		{"(emacs-pid)", "(emacs&-pid)"},
	}
	for _, test := range tests {
		input := test.unquoted
		want := test.quoted
		got := server_quote_arg(input)
		if got != want {
			t.Errorf("server_quote_arg(%q) = %q, want %q", input, got, want)
		}
	}
	for _, test := range tests {
		input := test.quoted
		want := test.unquoted
		got := server_unquote_arg(input)
		if got != want {
			t.Errorf("server_unquote_arg(%q) = %q, want %q", test.unquoted, got, test.quoted)
		}
	}
}
