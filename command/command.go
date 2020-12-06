package command

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Command is a frontend to exec.Cmd.
type Command struct {
	*exec.Cmd
}

// New configures a new command.
func New(line string) *Command {
	fields := strings.Fields(line)
	cmd := &Command{exec.Command(fields[0], fields[1:]...)}
	cmd.Env = os.Environ()
	return cmd
}

// OsPassthrough makes the command use the operating system standard IO.
func (c *Command) OsPassthrough() {
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
}

// Setenv sets an environment variable.
func (c *Command) Setenv(env ...string) {
	c.Env = append(c.Env, env...)
}

// WrapShell prepares the command for running inside a shell.
func (c *Command) WrapShell() {
	shell := New("/bin/sh")
	shell.Args = []string{"/bin/sh", "-c", fmt.Sprint(c)}
	shell.Stdin = c.Stdin
	shell.Stdout = c.Stdout
	shell.Stderr = c.Stderr
	shell.Env = c.Env
	*c = *shell
}

// RunPassthrough is a shorthand to run a command in passthrough mode.
func RunPassthrough(line string, env ...string) error {
	cmd := New(line)
	cmd.Setenv(env...)
	cmd.OsPassthrough()
	return cmd.Run()
}

// RunShell is a shorthand to run a command inside a shell.
func RunShell(line string, env ...string) error {
	cmd := New(line)
	cmd.Setenv(env...)
	cmd.WrapShell()
	return cmd.Run()
}
