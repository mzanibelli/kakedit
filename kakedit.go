// Package kakedit is a collection of tools to improve Kakoune session
// management.
package kakedit

import (
	"context"
	"fmt"
	"kakedit/internal/kakoune"
	"kakedit/internal/listener"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/sync/errgroup"
)

// Kakoune runs kak(1) within a persistent session.
func Kakoune(cwd string, args ...string) error {
	kak := kakoune.FromEnvironment()

	// If a target session is already set by environment, expect
	// it to be fully started.
	if kak.Session != "" {
		return kak.EditSession(args...).Run()
	}

	sess, err := os.ReadFile(filepath.Join(cwd, ".kaksession"))
	if err == nil {
		kak.Session = strings.TrimSpace(string(sess))
	} else if err != nil && !os.IsNotExist(err) {
		return err
	}

	// No session requested or it is already running.
	if kak.Session == "" || kak.Ping().Run() == nil {
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
		defer close(ping)
		for kak.Ping().Run() != nil {
			continue
		}
	}()

	select {
	case <-ping:
		return kak.EditSession(args...).Run()
	case <-time.After(1 * time.Second):
		return context.DeadlineExceeded
	}
}

// ExternalProgram runs an external program with a modified $EDITOR.
func ExternalProgram(pipe string, args ...string) error {
	kak, addr := kakoune.FromEnvironment(), listener.UniqueSocketPath()

	parent, cancel := context.WithCancel(context.Background())
	defer cancel()

	group, ctx := errgroup.WithContext(parent)

	group.Go(func() error {
		return listener.ListenContext(ctx, addr, listener.OnMessageFunc(func(data []byte) error {
			return kak.EditClientBulk(strings.Split(string(data), "\n"))
		}))
	})

	group.Go(func() error {
		defer cancel()

		cmd := exec.Command(args[0], args[1:]...)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		cmd.Env = append(
			os.Environ(),
			fmt.Sprintf("EDITOR=%s %s", pipe, addr),
			fmt.Sprintf("VISUAL=%s %s", pipe, addr),
		)

		return cmd.Run()
	})

	return group.Wait()
}
