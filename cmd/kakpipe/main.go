// kakpipe(1) writes a string to given socket.
package main

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
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

	// Kakoune supports jumping to a specific line number.
	// The jump is performed on the first file argument only.
	var files []string
	var lnum string
	for _, file := range args[2:] {
		if strings.HasPrefix(file, "+") {
			lnum = strings.TrimPrefix(file, "+")
			continue
		}

		abs, err := filepath.Abs(file)
		if err != nil {
			return err
		}

		files = append(files, abs)
	}

	// Add line number to first match.
	files[0] = files[0] + " " + lnum

	for _, file := range files {
		if _, err = fmt.Fprintln(conn, file); err != nil {
			return err
		}
	}

	return nil
}
