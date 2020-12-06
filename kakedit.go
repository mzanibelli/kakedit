package kakedit

import (
	"context"
	"fmt"
	"kakedit/client"
	"kakedit/command"
	"kakedit/kakoune"
	"kakedit/listener"
	"time"
)

// Holds the information to identify a kakoune session and its preferred
// client. A global variable is used to easily make this available
// inside the client or the server context.
var kak = kakoune.FromEnvironment()

// Pick runs the file picker and listens for edit requests.
func Pick(self, bin string, timeout time.Duration) error {
	var err error

	ctx, cancel := context.WithCancel(context.Background())
	lst, err := listener.ListenContext(ctx, timeout)
	if err != nil {
		return err
	}

	lst.Run(listener.OnMessageFunc(func(data []byte) error {
		return command.RunShell(kak.Edit(string(data)))
	}))

	cl := client.New(self, lst.Addr())

	err = command.RunPassthrough(bin,
		fmt.Sprintf("EDITOR=%s", cl), fmt.Sprintf("VISUAL=%s", cl))

	cancel()

	// Make sure we call Close() and catch any error if no previous error occured.
	if closeErr := lst.Close(); closeErr != ctx.Err() && err == nil {
		err = closeErr
	}

	return err
}

// Edit sends a given filename to the remote listener.
func Edit(bin, socket, file string) error { return client.New(bin, socket).Send(file) }
