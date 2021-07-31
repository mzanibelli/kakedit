// kakwrap(1) runs Kakoune inside a persistent session.
package main

import (
	"fmt"
	"kakedit"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	if err := run(os.Args...); err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", os.Args[0], err)
		os.Exit(1)
	}
}

func run(args ...string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	// Ignore interrupts like Kakoune does.
	signal.Ignore(syscall.SIGINT)

	return kakedit.Kakoune(cwd, args[1:]...)
}
