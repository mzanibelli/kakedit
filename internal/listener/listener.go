package listener

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"path"
	"time"
)

// UniqueSocketPath returns a path to a new socket in XDG_RUNTIME_DIR or in
// system temporary directory if unset.
func UniqueSocketPath() string {
	baseDir := os.Getenv("XDG_RUNTIME_DIR")
	if _, err := os.Stat(baseDir); os.IsNotExist(err) {
		baseDir = os.TempDir()
	}

	// Example output: /run/user/1000/kakedit.XXX according to pid.
	return path.Join(baseDir, fmt.Sprintf("kakedit.%d", os.Getpid()))
}

// Handler handles a new connection.
type Handler interface{ OnMessage(data []byte) error }

// OnMessageFunc is a function implementation of Handler.
type OnMessageFunc func([]byte) error

// OnMessage handles a new connection.
func (f OnMessageFunc) OnMessage(data []byte) error { return f(data) }

// ListenContext creates a stoppable listener.
func ListenContext(ctx context.Context, addr string, h Handler) error {
	lst, err := net.Listen("unix", addr)
	if err != nil {
		return err
	}
	defer lst.Close()

	unix := lst.(*net.UnixListener)

	const timeout = 20 * time.Millisecond
	for {
		if err := unix.SetDeadline(time.Now().Add(timeout)); err != nil {
			return err
		}

		conn, err := unix.Accept()
		select {
		case <-ctx.Done():
			return nil
		default:
			// handle connection
		}

		var ne net.Error
		if errors.As(err, &ne) && ne.Timeout() && ne.Temporary() {
			continue
		}

		if err != nil {
			return err
		}
		defer conn.Close()

		data, err := io.ReadAll(conn)
		if err != nil {
			return err
		}

		if err := h.OnMessage(data); err != nil {
			return err
		}
	}
}
