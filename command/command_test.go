package command_test

import (
	"kakedit/command"
	"os"
	"strings"
	"testing"
)

func TestCommand(t *testing.T) {
	cmd := command.New("foo bar")
	cmd.OsPassthrough()
	cmd.Setenv("bar=baz")
	cmd.WrapShell()

	if cmd.Stdout != os.Stdout {
		t.Error("output is not stdout")
	}
	if cmd.Stderr != os.Stderr {
		t.Error("error is not stderr")
	}
	if cmd.Stdin != os.Stdin {
		t.Error("input is not stdin")
	}

	assertEnv := func(env, name, value string) {
		parts := strings.Split(env, "=")
		if parts[0] == name && parts[1] != value {
			t.Errorf("want: %s, got: %s", value, parts[1])
		}
	}

	for _, env := range cmd.Env {
		assertEnv(env, "bar", "baz")
	}

	want := "/bin/sh -c foo bar"
	got := cmd.String()
	if cmd.String() != want {
		t.Errorf("want: %s, got: %s", want, got)
	}
}
