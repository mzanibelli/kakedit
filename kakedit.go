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
func Server(cmd string) error {
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
		cmd, err := kak.EditClient(string(data))
		if err != nil {
			return err
		}
		return command.RunShell(cmd)
	}))

	env := []string{
		fmt.Sprintf("EDITOR=%s -mode client %s", self, lst),
		fmt.Sprintf("VISUAL=%s -mode client %s", self, lst),
	}
	err = command.RunPassthrough(cmd, env...)

	cancel()

	// Make sure we call Close() and catch any error if no previous error occured.
	if closeErr := lst.Close(); closeErr != ctx.Err() && err == nil {
		err = closeErr
	}

	return err
}

// Local runs the file picker and replaces $EDITOR with a pre-connected
// Kakoune command.
func Local(cmd string) error {
	editor, err := kak.EditSession()
	if err != nil {
		return err
	}

	env := []string{
		fmt.Sprintf("EDITOR=%s", editor),
		fmt.Sprintf("VISUAL=%s", editor),
	}
	return command.RunPassthrough(cmd, env...)
}

// Client acts as a drop-in $EDITOR replacement and sends filenames to
// the server.
func Client(socket, file string) error { return listener.Send(socket, file) }
