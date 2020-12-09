// kakedit(1) runs an external tool with a modified $EDITOR.
package main

import (
	"fmt"
	"kakedit"
	"os"
	"os/exec"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		usage()
	}

	kakpipe, err := exec.LookPath("kakpipe")
	if err != nil {
		exit(err)
	}

	kakwrap, err := exec.LookPath("kakwrap")
	if err != nil {
		exit(err)
	}

	cmd := strings.Join(os.Args[1:], " ")

	if err := kakedit.ExternalProgram(cmd, kakpipe, kakwrap); err != nil {
		exit(err)
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s PROGRAM [ARGS...]\n", os.Args[0])
	os.Exit(1)
}

func exit(err error) {
	fmt.Fprintf(os.Stderr, "%s: %v", os.Args[0], err)
	os.Exit(1)
}
