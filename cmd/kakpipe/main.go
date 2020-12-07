// kakpipe(1) writes a string to given socket.
package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		usage()
	}

	addr := os.Args[1]
	file := strings.Join(os.Args[2:], " ")

	conn, err := net.Dial("unix", addr)
	if err != nil {
		exit(err)
	}

	defer conn.Close()

	if _, err = fmt.Fprint(conn, file); err != nil {
		exit(err)
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s ADDR FILE\n", os.Args[0])
	os.Exit(1)
}

func exit(err error) {
	fmt.Fprintf(os.Stderr, "%s: %v", os.Args[0], err)
	os.Exit(1)
}
