package picker_test

import (
	"bytes"
	"kakedit/picker"
	"os"
	"strings"
	"testing"
)

func TestPicker(t *testing.T) {
	pi := picker.New("foo")
	cmd := pi.Pick(bytes.NewBuffer([]byte("bar")))

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
		assertEnv(env, "EDITOR", "bar")
		assertEnv(env, "VISUAL", "bar")
	}
}
