package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	server_socket_dir = fmt.Sprintf("%s/emacs%d", os.TempDir(), os.Getuid())
	server_name       = "server"
)

var (
	eval = flag.String("e", "", "Eval an ELisp expression")
)

// server_quote_arg is alternative to Emacs's server-quote-arg
func server_quote_arg(arg string) string {
	r := regexp.MustCompile("[-&\n ]")
	return r.ReplaceAllStringFunc(arg, func(x string) string {
		switch x[0] {
		case '&':
			return "&&"
		case '-':
			return "&-"
		case '\n':
			return "&n"
		case ' ':
			return "&_"
		}
		return arg
	})
}

// server_unquote_arg is alternative to Emacs's server-unquote-arg
func server_unquote_arg(arg string) string {
	r := regexp.MustCompile("&.")
	return r.ReplaceAllStringFunc(arg, func(x string) string {
		switch x[1] {
		case '&':
			return "&"
		case '-':
			return "-"
		case '\n':
			return "\n"
		default:
			return " "
		}
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
	input := bufio.NewScanner(c)
	buf := new(bytes.Buffer)
	for input.Scan() {
		line := input.Text()
		var s string
		switch {
		case strings.HasPrefix(line, "-print "):
			s = line[len("-print "):]
		case strings.HasPrefix(line, "-print-nonl "):
			s = line[len("-print-nonl "):]
		case strings.HasPrefix(line, "-error "):
			s = line[len("-error "):]
		default:				// such as -emacs-pid 274
			continue
		}
		buf.WriteString(server_unquote_arg(s))
	}
	if buf.Len() > 0 {
		buf.WriteByte('\n')
		if _, err := io.Copy(w, buf); err != nil {
			return err
		}
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
