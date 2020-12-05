package kakoune

import (
	"fmt"
	"os"
	"os/exec"
)

// Instance holds information to make a connection to an existing
// Kakoune client.
type Instance struct {
	session string
	client  string
}

// New returns a new Kakoune instance.
func New(session, client string) *Instance {
	return &Instance{session, client}
}

// Edit sends an edit command to the instance.
func (i Instance) Edit(file string) *exec.Cmd {
	cmd := fmt.Sprintf("evaluate-commands -verbatim -client %s edit \"%s\"",
		i.client, file)
	return i.makeShellCmd(cmd)
}

func (i Instance) makeShellCmd(str string) *exec.Cmd {
	shellCmd := fmt.Sprintf("echo '%s' | kak -p %s", str, i.session)
	cmd := exec.Command("/bin/sh", "-c", shellCmd)
	cmd.Env = os.Environ()
	return cmd
}
