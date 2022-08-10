package kakoune

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Kakoune holds information to make a connection to an existing
// Kakoune client.
type Kakoune struct {
	Session string
	Client  string
}

// FromEnvironment returns a new Kakoune instance.
func FromEnvironment() *Kakoune {
	return &Kakoune{
		os.Getenv("kak_session"),
		os.Getenv("kak_client"),
	}
}

// EditClient sends an edit command to an existing Kakoune client.
func (kak *Kakoune) EditClient(file string) *exec.Cmd {
	shell := fmt.Sprintf("echo 'evaluate-commands -verbatim -client %s edit -existing %s' | %s -p %s",
		kak.Client, strings.TrimSpace(file), "kak", kak.Session)
	return exec.Command("/bin/sh", "-c", shell)
}

// Ping succeeds if the session is up and running.
func (kak *Kakoune) Ping() *exec.Cmd {
	shell := fmt.Sprintf("echo nop | kak -p %s", kak.Session)
	return exec.Command("/bin/sh", "-c", shell)
}

// EditClientBulk sends edit commands for each received file.
func (kak *Kakoune) EditClientBulk(files []string) error {
	for _, file := range files {
		if file == "" {
			continue
		}
		if err := kak.EditClient(file).Run(); err != nil {
			return err
		}
	}
	return nil
}

// EditSession starts a new client connected to the session.
func (kak *Kakoune) EditSession(args ...string) *exec.Cmd {
	session := make([]string, 0, 2)
	if kak.Session != "" {
		session = append(session, "-c", kak.Session)
	}

	//nolint:gosec
	cmd := exec.Command("kak", append(session, args...)...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd
}

// StartSession starts a new session in a given directory.
// TODO: implement cross-platform setsid behavior.
func (kak *Kakoune) StartSession(cwd string) *exec.Cmd {
	//nolint:gosec
	return exec.Command("kak", "-s", kak.Session, "-d", "-E", fmt.Sprintf("cd '%s'", cwd))
}
