package kakoune

import (
	"fmt"
	"os"
)

const (
	editClient  = "echo 'evaluate-commands -verbatim -client %s edit \"%s\"' | kak -p %s"
	editSession = "kak -c %s"
)

// Kakoune holds information to make a connection to an existing
// Kakoune client.
type Kakoune struct {
	session string
	client  string
}

// FromEnvironment returns a new Kakoune instance.
func FromEnvironment() Kakoune {
	return Kakoune{os.Getenv("kak_session"), os.Getenv("kak_client")}
}

// EditClient sends an edit command to an existing Kakoune client.
func (kak Kakoune) EditClient(file string) string {
	return fmt.Sprintf(editClient, kak.client, file, kak.session)
}

// EditSession starts a new client connected to the same session.
func (kak Kakoune) EditSession() string {
	return fmt.Sprintf(editSession, kak.session)
}
