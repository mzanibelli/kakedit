package kakoune_test

import (
	"kakedit/kakoune"
	"reflect"
	"testing"
)

func TestEditCommandFormat(t *testing.T) {
	kak := kakoune.New("foo", "bar")
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
