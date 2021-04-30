package listener_test

import (
	"context"
	"errors"
	"fmt"
	"kakedit/internal/listener"
	"net"
	"os"
	"sync"
	"testing"
)

func TestListener(t *testing.T) {
	mess := "hello"

	ctx, cancel := context.WithCancel(context.Background())
	lst, err := listener.ListenContext(ctx)
	if err != nil {
		t.Errorf("could not initialize listener: %v", err)
	}

	wg := new(sync.WaitGroup)
	wg.Add(1)

	lst.HandleFunc(func(data []byte) error {
		if string(data) != mess {
			t.Errorf("want: %s, got: %s", mess, data)
		}
		wg.Done()
		return nil
	})

	conn, err := net.Dial("unix", lst.Addr().String())
	if err != nil {
		t.Errorf("could not connect to listener: %v", err)
	}
	_, err = fmt.Fprint(conn, mess)
	if err != nil {
		t.Errorf("could not send message: %v", err)
	}
	conn.Close()

	wg.Wait()

	cancel()

	if err := lst.Close(); err != nil && !errors.Is(err, ctx.Err()) {
		t.Error(err)
	}
	if _, err := os.Stat(lst.Addr().String()); os.IsExist(err) {
		t.Error("socket was not removed on close")
	}
}
