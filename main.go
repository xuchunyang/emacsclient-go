package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"io"
	"flag"
	"regexp"
)

var (
	server_socket_dir = fmt.Sprintf("%s/emacs%d", os.TempDir(), os.Getuid())
	server_name       = "server"
)

var (
	eval = flag.String("e", "", "Eval an ELisp expression")
)

// server_quote_arg is ported from server-quote-arg
func server_quote_arg (arg string) string {
	r := regexp.MustCompile("[-&\n ]")
	return r.ReplaceAllStringFunc(arg, func(x string) string {
		switch x[0] {
		case '&':
			return "&&"
		case '-':
			return "&_"
		case '\n':
			return "&n"
		case ' ':
			return "&_"
		}
		return arg
	})
}

func ensureTrailingNewline(s string) string {
	if s == "" {
		return s
	}
	if s[len(s)-1] == '\n' {
		return s
	}
	return s + "\n"
}

// Eval evaluates an ELisp expression
func Eval(c net.Conn, w io.Writer, expr string) error {
	command := "-eval " + ensureTrailingNewline(server_quote_arg(expr))
	if _, err := io.WriteString(c, command); err != nil {
		return err
	}
	if _, err := io.Copy(w, c); err != nil {
		return err
	}
	return nil
}

func main() {
	flag.Parse()
	if *eval == "" {
		flag.Usage()
		os.Exit(1)
	}

	server_file := filepath.Join(server_socket_dir, server_name)
	conn, err := net.Dial("unix", server_file)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	if err := Eval(conn, os.Stdout, *eval); err != nil {
		log.Fatal(err)
	}
}
