package kakoune

import (
	"errors"
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
func (kak Kakoune) EditClient(file string) (string, error) {
	if kak.session == "" || kak.client == "" {
		return "", errors.New("missing required environment variable: kak_session, kak_client")
	}
	return fmt.Sprintf(editClient, kak.client, file, kak.session), nil
}

// EditSession starts a new client connected to the same session.
func (kak Kakoune) EditSession() (string, error) {
	if kak.session == "" {
		return "", errors.New("missing required environment variable: kak_session")
	}
	return fmt.Sprintf(editSession, kak.session), nil
}
