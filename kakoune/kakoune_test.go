package kakoune_test

import (
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

	var want, got string

	want = strings.Join([]string{
		"echo",
		"'evaluate-commands -verbatim -client bar edit \"test.txt\"'",
		"| kak -p foo",
	}, " ")
	got = kak.EditClient("test.txt")
	if want != got {
		t.Errorf("want: %s, got: %s", want, got)
	}

	want = "kak -c foo"
	got = kak.EditSession()
	if want != got {
		t.Errorf("want: %s, got: %s", want, got)
	}
}
