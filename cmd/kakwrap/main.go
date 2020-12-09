// kakwrap(1) runs Kakoune inside a persistent session.
package main

import (
	"fmt"
	"kakedit"
	"os"
)

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		exit(err)
	}

	if err := kakedit.Kakoune(cwd, os.Args[1:]...); err != nil {
		exit(err)
	}
}

func exit(err error) {
	fmt.Fprintf(os.Stderr, "%s: %v", os.Args[0], err)
	os.Exit(1)
}
