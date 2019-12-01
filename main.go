package main

import (
	"bufio"
	"log"
	"bytes"
	"flag"
	"fmt"
	"io"
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
	evalFlag = flag.Bool("e", false, "Evaluate the FILE arguments as ELisp expressions")
	nowaitFlag = flag.Bool("n", false, "Don't wait for the server to return")
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

func process(c net.Conn, output io.Writer, command string) error {
	command = ensureTrailingNewline(command)
	if _, err := io.WriteString(c, command); err != nil {
		return err
	}
	input := bufio.NewScanner(c)
	buf := new(bytes.Buffer)
	var newline bool
	for input.Scan() {
		line := input.Text()
		var s string
		switch {
		case strings.HasPrefix(line, "-emacs-pid"):
			continue
		case strings.HasPrefix(line, "-print "):
			s = line[len("-print "):]+"\n"
			newline = true
		case strings.HasPrefix(line, "-print-nonl "):
			if newline {
				buf.Truncate(buf.Len()-1)
			}
			s = line[len("-print-nonl "):]
			newline = false
		case strings.HasPrefix(line, "-error "):
			s = line[len("-error "):]+"\n"
			newline = true
		default:
			log.Printf("%q is not supported\n", line)
			continue
		}
		buf.WriteString(server_unquote_arg(s))
	}
	if buf.Len() > 0 {
		if _, err := io.Copy(output, buf); err != nil {
			return err
		}
	}
	return nil
}

func connect() (net.Conn, error) {
	server_file := filepath.Join(server_socket_dir, server_name)
	return net.Dial("unix", server_file)
}

func buildCommand() string {
	var commands []string
	if *nowaitFlag {
		commands = append(commands, "-nowait")
	}
	for _, arg := range flag.Args() {
		var cmd string
		if *evalFlag {
			cmd = "-eval " + server_quote_arg(arg)
		} else {
			cmd = "-file " + server_quote_arg(arg)
		}
		commands = append(commands, cmd)
	}
	return strings.Join(commands, " ")
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS] FILE...\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()
	if flag.NArg() == 0 {
		fmt.Fprintf(os.Stderr, "%s: file name or argument required\n", os.Args[0])
		flag.Usage()
		os.Exit(1)
	}
	conn, err := connect()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s", os.Args[0], err)
		os.Exit(1)
	}
	defer conn.Close()
	if err := process(conn, os.Stdout, buildCommand()); err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s", os.Args[0], err)
		os.Exit(1)
	}
}
