// kakpipe(1) writes a string to given socket.
package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		usage()
	}

	conn, err := net.Dial("unix", os.Args[1])
	if err != nil {
		exit(err)
	}

	defer conn.Close()

	for _, file := range os.Args[2:] {
		if _, err = fmt.Fprintln(conn, file); err != nil {
			exit(err)
		}
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
