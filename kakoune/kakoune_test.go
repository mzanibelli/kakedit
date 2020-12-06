package kakoune_test

import (
	"fmt"
	"kakedit/kakoune"
	"os"
	"strings"
	"testing"
)

func TestEditCommandFormat(t *testing.T) {
	if err := os.Setenv("kak_session", "foo"); err != nil {
		t.Error(err)
	}
	if err := os.Setenv("kak_client", "bar"); err != nil {
		t.Error(err)
	}

	kak := kakoune.FromEnvironment()
	cmd := kak.Edit("test.txt")

	want := strings.Join([]string{
		"echo",
		"'evaluate-commands -verbatim -client bar edit \"test.txt\"'",
		"| kak -p foo",
	}, " ")

	if fmt.Sprint(cmd) != want {
		t.Errorf("want: %s\ngot: %s", want, fmt.Sprint(cmd))
	}
}
