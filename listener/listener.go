package listener

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path"
	"time"
)

var (
	// ErrCaughtStopSignal is returned when the listener purposefully
	// stopped listening for incoming connections.
	ErrCaughtStopSignal = errors.New("caught stop signal")
)

// Listener is a stoppable UnixListener.
type Listener struct {
	*net.UnixListener
	addr    string
	timeout time.Duration
	done    chan struct{}
}

// Listen creates a stoppable listener.
func Listen(timeout time.Duration) (*Listener, error) {
	addr := makeUniqueSocketPath()
	unix, err := net.Listen("unix", addr)
	if err != nil {
		return nil, err
	}
	return &Listener{
		unix.(*net.UnixListener),
		addr,
		timeout,
		make(chan struct{}),
	}, nil
}

// Close stops waiting for the next incoming connection.
func (l *Listener) Close() {
	defer l.UnixListener.Close()
	defer close(l.done)
	l.done <- struct{}{}
}

// Addr returns the path to the socket being listened to.
func (l Listener) Addr() string { return l.addr }

// Handler handles a new connection.
type Handler interface{ OnMessage(data []byte) error }

// OnMessageFunc is a function implementation of Handler.
type OnMessageFunc func([]byte) error

// OnMessage handles a new connection.
func (f OnMessageFunc) OnMessage(data []byte) error { return f(data) }

// Run starts waiting for an incoming connection in a separate goroutine.
func (l *Listener) Run(handler Handler) <-chan error {
	ec := make(chan error, 1)
	go func() {
		defer close(ec)
		ec <- l.run(handler)
	}()
	return ec
}

func (l *Listener) run(handler Handler) error {
	for {
		l.SetDeadline(time.Now().Add(l.timeout))
		conn, err := l.UnixListener.Accept()
		select {
		case <-l.done:
			return ErrCaughtStopSignal
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
		data, err := ioutil.ReadAll(conn)
		if err != nil {
			return err
		}
		if err := handler.OnMessage(data); err != nil {
			return err
		}
	}
}

func makeUniqueSocketPath() string {
	baseDir := os.Getenv("XDG_RUNTIME_DIR")
	if _, err := os.Stat(baseDir); os.IsNotExist(err) {
		baseDir = os.TempDir()
	}
	return path.Join(baseDir,
		fmt.Sprintf("kakedit.%d", os.Getpid()))
}
