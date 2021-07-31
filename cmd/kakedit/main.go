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
	if err := run(os.Args...); err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", os.Args[0], err)
		os.Exit(1)
	}
}

func run(args ...string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: %s PROGRAM [ARGS...]", args[0])
	}

	kakpipe, err := exec.LookPath("kakpipe")
	if err != nil {
		return err
	}

	cmd := strings.Join(args[1:], " ")

	return kakedit.ExternalProgram(cmd, kakpipe)
}
