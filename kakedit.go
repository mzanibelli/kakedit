// Package kakedit is a wrapper to an external program. It makes $EDITOR
// invocations connect to an existing Kakoune instance.
package kakedit

import (
	"context"
	"fmt"
	"kakedit/internal/kakoune"
	"kakedit/internal/listener"
	"os"
	"os/exec"
	"time"
)

// Run runs the server and waits for edit requests.
func Run(cmd, pipe string) error {
	kak := kakoune.FromEnvironment()

	if kak.UnknownRemote() {
		return runShell(cmd, kak.EditSession())
	}

	var err error

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	lst, err := listener.ListenContext(ctx, 20*time.Millisecond)
	if err != nil {
		return err
	}

	lst.Run(listener.OnMessageFunc(func(data []byte) error {
		return runShell(kak.EditClient(string(data)), "")
	}))

	editor := fmt.Sprintf("%s %s", pipe, lst.Addr())

	err = runShell(cmd, editor)

	cancel()

	// Make sure we call Close() and catch any error if no previous error occured.
	if closeErr := lst.Close(); closeErr != ctx.Err() && err == nil {
		err = closeErr
	}

	return err
}

// Run a command inside a shell. Replace environment variables EDITOR
// and VISUAL with the given value.
func runShell(line string, editor string) error {
	cmd := exec.Command("/bin/sh", "-c", line)

	// If no $EDITOR given, run the command without system IO.
	if editor == "" {
		cmd.Stdin = nil
		cmd.Stdout = nil
		cmd.Stderr = nil
	} else {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	cmd.Env = append(os.Environ(),
		fmt.Sprintf("EDITOR=%s", editor), fmt.Sprintf("VISUAL=%s", editor))

	return cmd.Run()
}
