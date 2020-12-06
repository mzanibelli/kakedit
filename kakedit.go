package kakedit

import (
	"context"
	"fmt"
	"kakedit/kakoune"
	"kakedit/listener"
	"os"
	"os/exec"
	"time"
)

// Holds the information to identify a kakoune session and its preferred
// client. A global variable is used to easily make this available
// inside the client or the server context.
var kak = kakoune.FromEnvironment()

// DefaultTimeout is the maximum delay to shutdown the listener.
const DefaultTimeout time.Duration = 20 * time.Millisecond

// Server runs the file picker and waits for edit requests. Requests are
// then forwarded to an existing Kakoune client.
func Server(cmd string) error {
	if kak.UnknownRemote() {
		return Local(cmd)
	}

	var err error

	self, err := os.Executable()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	lst, err := listener.ListenContext(ctx, DefaultTimeout)
	if err != nil {
		return err
	}

	lst.Run(listener.OnMessageFunc(func(data []byte) error {
		return runShell(kak.EditClient(string(data)), "")
	}))

	// Force client to run in silent mode for now. One can change
	// the source for debugging purposes.
	err = runShell(cmd, fmt.Sprintf("%s -silent -mode client %s", self, lst))

	cancel()

	// Make sure we call Close() and catch any error if no previous error occured.
	if closeErr := lst.Close(); closeErr != ctx.Err() && err == nil {
		err = closeErr
	}

	return err
}

// Local runs the file picker and replaces $EDITOR with a pre-connected
// Kakoune command. This fallbacks to a brand new Kakoune session if
// the environment variable is empty.
func Local(cmd string) error { return runShell(cmd, kak.EditSession()) }

// Client acts as a drop-in $EDITOR replacement and sends filenames to
// the server.
func Client(socket, file string) error { return listener.Send(socket, file) }

// Run a command inside a shell using system IO. Replace environment
// variables EDITOR and VISUAL with the given value.
func runShell(line string, editor string) error {
	cmd := exec.Command("/bin/sh", "-c", line)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("EDITOR=%s", editor), fmt.Sprintf("VISUAL=%s", editor))
	return cmd.Run()
}
