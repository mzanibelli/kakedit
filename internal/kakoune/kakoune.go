package kakoune

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"strings"
)

// Kakoune holds information to make a connection to an existing
// Kakoune client.
type Kakoune struct {
	Bin     string
	Session string
	Client  string
}

// FromEnvironment returns a new Kakoune instance.
func FromEnvironment() *Kakoune {
	bin, err := exec.LookPath("kak")

	// TODO: do not use panic()
	if err != nil {
		panic("kak(1) is missing in $PATH")
	}

	return &Kakoune{
		bin,
		os.Getenv("kak_session"),
		os.Getenv("kak_client"),
	}
}

// UnknownRemote returns true if the current environment does not allow
// targeting a remote Kakoune instance.
func (kak *Kakoune) UnknownRemote() bool {
	return kak.Session == "" || kak.Client == ""
}

// EditClient sends an edit command to an existing Kakoune client.
func (kak *Kakoune) EditClient(file string) *exec.Cmd {
	shell := fmt.Sprintf("echo 'evaluate-commands -verbatim -client %s edit -existing \"%s\"' | %s -p %s",
		kak.Client, strings.TrimSpace(file), kak.Bin, kak.Session)
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

	cmd := exec.Command(kak.Bin, append(session, args...)...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd
}

// StartSession starts a new session in a given directory.
func (kak *Kakoune) StartSession(cwd string) *exec.Cmd {
	return exec.Command(kak.Bin, "-s", kak.Session, "-d", "-E", fmt.Sprintf("cd '%s'", cwd))
}

// SetUniqueSessionName sets an unique session name for a given path.
func (kak *Kakoune) SetUniqueSessionName(cwd string) {
	h := md5.New()
	io.WriteString(h, path.Base(cwd))
	kak.Session = fmt.Sprintf("%x", h.Sum(nil))
}
