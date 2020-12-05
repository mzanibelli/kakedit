package kakedit

import (
	"errors"
	"kakedit/editor"
	"kakedit/kakoune"
	"kakedit/listener"
	"kakedit/picker"
	"sync"
	"time"
)

// Pick runs the file picker and listens for edit requests.
func Pick(edBin, piBin, session, client string, timeout time.Duration) error {
	lst, err := listener.Listen(timeout)
	if err != nil {
		return err
	}

	var once sync.Once
	defer once.Do(lst.Close)

	kak := kakoune.New(session, client)

	ec := lst.Run(listener.OnMessageFunc(func(data []byte) error {
		return kak.Edit(string(data)).Run()
	}))

	ed := editor.New(edBin, lst.Addr())
	pi := picker.New(piBin)

	if err := pi.Pick(ed).Run(); err != nil {
		return err
	}

	once.Do(lst.Close)

	if err := <-ec; err != nil &&
		!errors.Is(err, listener.ErrCaughtStopSignal) {
		return err
	}

	return nil
}

// Edit sends a given filename to the remote listener.
func Edit(bin, socket, file string) error { return editor.New(bin, socket).Edit(file) }
