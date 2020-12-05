package editor

import (
	"fmt"
	"net"
)

// Editor edits files.
type Editor struct {
	bin  string
	addr string
}

func (e Editor) String() string {
	return fmt.Sprintf("%s %s", e.bin, e.addr)
}

// New creates an editor.
func New(bin, addr string) *Editor {
	return &Editor{bin, addr}
}

// Edit sends a remote edit request containing the target file.
func (e Editor) Edit(file string) error {
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
