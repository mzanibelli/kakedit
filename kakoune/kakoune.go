package kakoune

import (
	"fmt"
	"os"
)

const (
	editCmd = "echo 'evaluate-commands -verbatim -client %s edit \"%s\"' | kak -p %s"
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
func (i Kakoune) Edit(file string) string {
	return fmt.Sprintf(editCmd, i.client, file, i.session)
}
