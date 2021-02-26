// Package kakedit is a collection of tools to improve Kakoune session
// management.
package kakedit

import (
	"context"
	"errors"
	"fmt"
	"kakedit/internal/kakoune"
	"kakedit/internal/listener"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"
)

// Kakoune runs kak(1) within a persistent session.
func Kakoune(cwd string, args ...string) error {
	kak := kakoune.FromEnvironment()

	// If a target session is already set by environment, expected
	// it to be fully started.
	if kak.Session != "" {
		return kak.EditSession(args...).Run()
	}

	sess, err := os.ReadFile(path.Join(cwd, ".kaksession"))
	if err == nil {
		kak.Session = strings.TrimSpace(string(sess))
	} else if err != nil && !os.IsNotExist(err) {
		return err
	}

	// Is a specific session requested?
	if kak.Session == "" {
		return kak.EditSession(args...).Run()
	}

	// Did somebody else start the session?
	if err := kak.Ping().Run(); err == nil {
		return kak.EditSession(args...).Run()
	}

	// Do not check for actual errors on background session start. Let
	// it silently fail if address already in use.
	if err := kak.StartSession(cwd).Start(); err != nil {
		return err
	}

	// Wait for the session to be truly started.
	// See: https://github.com/mawww/kakoune/issues/3618
	ping := make(chan struct{})
	go func() {
		for {
			if err := kak.Ping().Run(); err == nil {
				close(ping)
				break
			}
		}
	}()

	select {
	case <-ping:
		return kak.EditSession(args...).Run()
	case <-time.After(1 * time.Second):
		return errors.New("timeout waiting for session")
	}
}

// ExternalProgram runs an external program with a modified $EDITOR.
func ExternalProgram(shell, kakpipe, kakwrap string) error {
	kak := kakoune.FromEnvironment()

	// If we cannot connect to a running client, trust kakwrap(1)
	// to nicely create a new one.
	if kak.UnknownRemote() {
		cmd := exec.Command(kakwrap)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}

	var err error

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	lst, err := listener.ListenContext(ctx, 20*time.Millisecond)
	if err != nil {
		return err
	}

	lst.Run(listener.OnMessageFunc(func(data []byte) error {
		return kak.EditClientBulk(strings.Split(string(data), "\n"))
	}))

	// Run inside a shell to allow tricks like `$EDITOR $(fzf)`.
	cmd := exec.Command("/bin/sh", "-c", shell)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Replace $EDITOR with kakpipe(1) pre-connected to the socket.
	cmd.Env = append(
		os.Environ(),
		fmt.Sprintf("EDITOR=%s %s", kakpipe, lst.Addr()),
		fmt.Sprintf("VISUAL=%s %s", kakpipe, lst.Addr()),
	)

	err = cmd.Run()

	cancel()

	// Make sure we call Close() and catch any error if no previous error occured.
	if closeErr := lst.Close(); closeErr != ctx.Err() && err == nil {
		err = closeErr
	}

	return err
}
