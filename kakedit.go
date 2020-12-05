package kakedit

import (
	"context"
	"kakedit/editor"
	"kakedit/kakoune"
	"kakedit/listener"
	"kakedit/picker"
	"time"
)

// Pick runs the file picker and listens for edit requests.
func Pick(edBin, piBin, session, client string, timeout time.Duration) error {
	var err error

	kak := kakoune.New(session, client)

	ctx, cancel := context.WithCancel(context.Background())
	lst, err := listener.ListenContext(ctx, timeout)
	if err != nil {
		return err
	}
	lst.Run(listener.OnMessageFunc(func(data []byte) error {
		return kak.Edit(string(data)).Run()
	}))

	ed := editor.New(edBin, lst.Addr())
	pi := picker.New(piBin)

	err = pi.Pick(ed).Run()

	cancel()

	// Make sure we call Close() and catch any error if no previous error occured.
	if closeErr := lst.Close(); closeErr != ctx.Err() && err == nil {
		err = closeErr
	}

	return err
}

// Edit sends a given filename to the remote listener.
func Edit(bin, socket, file string) error { return editor.New(bin, socket).Edit(file) }
