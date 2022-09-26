package listener_test

import (
	"context"
	"fmt"
	"kakedit/internal/listener"
	"net"
	"os"
	"sync"
	"testing"
	"time"
)

func TestListener(t *testing.T) {
	mess, addr := "hello", listener.UniqueSocketPath()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg := new(sync.WaitGroup)
	wg.Add(1)

	h := listener.OnMessageFunc(func(data []byte) error {
		defer wg.Done()
		if string(data) != mess {
			t.Errorf("want: %s, got: %s", mess, data)
		}
		return nil
	})

	go func() {
		if err := listener.ListenContext(ctx, addr, h); err != nil {
			t.Error(err)
		}
	}()

	timeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := send(timeout, addr, mess); err != nil {
		t.Error(err)
	}

	wg.Wait()

	cancel()

	if _, err := os.Stat(addr); os.IsExist(err) {
		t.Error("socket was not removed on close")
	}
}

func send(ctx context.Context, addr, mess string) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			// pass
		}

		conn, err := net.Dial("unix", addr)
		if err != nil {
			continue
		}
		defer conn.Close()

		if _, err = fmt.Fprint(conn, mess); err != nil {
			return err
		}

		return nil
	}
}
