// kakpipe(1) writes a string to given socket.
package main

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
)

func main() {
	if err := run(os.Args...); err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", os.Args[0], err)
		os.Exit(1)
	}
}

func run(args ...string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: %s ADDR FILE", args[0])
	}

	conn, err := net.Dial("unix", args[1])
	if err != nil {
		return err
	}

	defer conn.Close()

	for _, file := range args[2:] {
		abs, err := filepath.Abs(file)
		if err != nil {
			return err
		}
		if _, err = fmt.Fprintln(conn, abs); err != nil {
			return err
		}
	}

	return nil
}
