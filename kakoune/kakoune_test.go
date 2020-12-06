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
	got, err := kak.EditClient("test.txt")
	if err != nil {
		t.Error(err)
	}
	if want != got {
		t.Errorf("want: %s, got: %s", want, got)
	}

	want = "kak -c foo"
	got, err = kak.EditSession()
	if err != nil {
		t.Error(err)
	}
	if want != got {
		t.Errorf("want: %s, got: %s", want, got)
	}
}

func TestFailureOnMissingEnvironment(t *testing.T) {
	if err := os.Setenv("kak_session", ""); err != nil {
		t.Error(err)
	}
	if err := os.Setenv("kak_client", ""); err != nil {
		t.Error(err)
	}

	kak := kakoune.FromEnvironment()

	if _, err := kak.EditClient("foo.txt"); err == nil {
		t.Error("should not work without both kak_session and kak_client specified")
	}
	if _, err := kak.EditSession(); err == nil {
		t.Error("should not work without kak_session specified")
	}
}
