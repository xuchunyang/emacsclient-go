# emacsclient-go - An Emacsclient in Go

emacsclient-go is a client for the [(emacs) Emacs
Server](https://www.gnu.org/software/emacs/manual/html_node/emacs/Emacs-Server.html)
written in Go. No one will need it because of the official `emacsclient(1)`. I
wrote it to learn how Emacs Server works and practice the Go language.

It supports a part of `emacsclient(1)`'s options:

    ~ $ emacsclient-go -h
    Usage: emacsclient-go [OPTIONS] FILE...
    Tell the Emacs server to visit the specified files.
    Every FILE can be either just a FILENAME or [+LINE[:COLUMN]] FILENAME.
    
    The following OPTIONS are accepted:
      -c    Create a new frame instead of trying to use the current Emacs frame
      -e    Evaluate the FILE arguments as ELisp expressions
      -n    Don't wait for the server to return
      -s string
            Set filename of the UNIX socket for communication (default "server")
      -u    Don't display return values from the server
