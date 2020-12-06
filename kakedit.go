package kakedit

import (
	"context"
	"fmt"
	"kakedit/command"
	"kakedit/kakoune"
	"kakedit/listener"
	"os"
	"time"
)

// Holds the information to identify a kakoune session and its preferred
// client. A global variable is used to easily make this available
// inside the client or the server context.
var kak = kakoune.FromEnvironment()

// DefaultTimeout is the maximum delay to shutdown the listener.
const DefaultTimeout time.Duration = 200 * time.Millisecond

// Server runs the file picker and waits for edit requests. Requests are
// then forwarded to an existing Kakoune client.
func Server(bin string) error {
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
		return command.RunShell(kak.EditClient(string(data)))
	}))

	env := []string{
		fmt.Sprintf("EDITOR=%s %s", self, lst),
		fmt.Sprintf("VISUAL=%s %s", self, lst),
	}
	err = command.RunPassthrough(bin, env...)

	cancel()

	// Make sure we call Close() and catch any error if no previous error occured.
	if closeErr := lst.Close(); closeErr != ctx.Err() && err == nil {
		err = closeErr
	}

	return err
}

// Client acts as a drop-in $EDITOR replacement and sends filenames to
// the server.
func Client(socket, file string) error { return listener.Send(socket, file) }

// Local runs the file picker and replaces $EDITOR with a pre-connected
// Kakoune command.
func Local(bin string) error {
	env := []string{
		fmt.Sprintf("EDITOR=%s", kak.EditSession()),
		fmt.Sprintf("VISUAL=%s", kak.EditSession()),
	}
	return command.RunPassthrough(bin, env...)
}
