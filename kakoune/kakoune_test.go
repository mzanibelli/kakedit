package kakoune_test

import (
	"kakedit/kakoune"
	"os"
	"reflect"
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

	want := []string{
		"/bin/sh",
		"-c",
		"echo 'evaluate-commands -verbatim -client bar edit \"test.txt\"' | kak -p foo",
	}

	if !reflect.DeepEqual(cmd.Args, want) {
		t.Errorf("want: %v\ngot: %v", want, cmd.Args)
	}
}
