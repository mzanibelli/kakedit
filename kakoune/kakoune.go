package kakoune

import (
	"fmt"
	"os"
	"os/exec"
)

// Kakoune holds information to make a connection to an existing
// Kakoune client.
type Kakoune struct {
	session string
	client  string
}

// FromEnvironment returns a new Kakoune instance.
func FromEnvironment() *Kakoune {
	return &Kakoune{os.Getenv("kak_session"), os.Getenv("kak_client")}
}

// Edit sends an edit command to the instance.
func (i Kakoune) Edit(file string) *exec.Cmd {
	cmd := fmt.Sprintf("evaluate-commands -verbatim -client %s edit \"%s\"",
		i.client, file)
	return i.makeShellCmd(cmd)
}

func (i Kakoune) makeShellCmd(str string) *exec.Cmd {
	shellCmd := fmt.Sprintf("echo '%s' | kak -p %s", str, i.session)
	cmd := exec.Command("/bin/sh", "-c", shellCmd)
	cmd.Env = os.Environ()
	return cmd
}
