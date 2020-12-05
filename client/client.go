package client

import (
	"fmt"
	"net"
)

// Client passes filenames to the server. It is used as a drop-in
// replacement for $EDITOR.
type Client struct {
	bin  string
	addr string
}

func (e Client) String() string {
	return fmt.Sprintf("%s %s", e.bin, e.addr)
}

// New creates a client.
func New(bin, addr string) *Client {
	return &Client{bin, addr}
}

// Send sends a remote edit request containing the target file.
func (e Client) Send(file string) error {
	conn, err := net.Dial("unix", e.addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = fmt.Fprint(conn, file)
	if err != nil {
		return err
	}

	return nil
}
