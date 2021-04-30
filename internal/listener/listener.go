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

// Listener is a stoppable UnixListener.
type Listener struct {
	*net.UnixListener
	ctx context.Context
	err chan error
}

func makeUniqueSocketPath() string {
	baseDir := os.Getenv("XDG_RUNTIME_DIR")
	if _, err := os.Stat(baseDir); os.IsNotExist(err) {
		baseDir = os.TempDir()
	}

	// Example output: /run/user/1000/kakedit.XXX according to pid.
	return path.Join(baseDir, fmt.Sprintf("kakedit.%d", os.Getpid()))
}

// ListenContext creates a stoppable listener.
func ListenContext(ctx context.Context) (*Listener, error) {
	unix, err := net.Listen("unix", makeUniqueSocketPath())
	if err != nil {
		return nil, err
	}

	return &Listener{
		unix.(*net.UnixListener),
		ctx,
		make(chan error),
	}, nil
}

// Close stops waiting for the next incoming connection.
func (l *Listener) Close() error {
	defer l.UnixListener.Close()
	defer close(l.err)
	return <-l.err
}

// Handler handles a new connection.
type Handler interface{ OnMessage(data []byte) error }

// OnMessageFunc is a function implementation of Handler.
type OnMessageFunc func([]byte) error

// OnMessage handles a new connection.
func (f OnMessageFunc) OnMessage(data []byte) error { return f(data) }

// Handle starts waiting for an incoming connection in a separate goroutine.
func (l *Listener) Handle(handler Handler) { go func() { l.err <- l.handle(handler) }() }

// HandleFunc allows to pass an anonymous function to Handle.
func (l *Listener) HandleFunc(f func([]byte) error) { l.Handle(OnMessageFunc(f)) }

func (l *Listener) handle(handler Handler) error {
	const timeout = 20 * time.Millisecond

	for {
		if err := l.SetDeadline(time.Now().Add(timeout)); err != nil {
			return err
		}

		conn, err := l.UnixListener.Accept()
		select {
		case <-l.ctx.Done():
			return l.ctx.Err()
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

		if err := handler.OnMessage(data); err != nil {
			return err
		}
	}
}
