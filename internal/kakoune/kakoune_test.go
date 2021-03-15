package kakoune_test

import (
	"kakedit/internal/kakoune"
	"os"
	"testing"
)

func TestKakoune(t *testing.T) {
	t.Run("it should produce a correct remote edit command", func(t *testing.T) {
		os.Setenv("kak_session", "foo")
		os.Setenv("kak_client", "bar")

		SUT := kakoune.FromEnvironment()

		want := `/bin/sh -c echo 'evaluate-commands -verbatim -client bar edit -existing "a.txt"' | kak -p foo`
		got := SUT.EditClient("a.txt").String()

		if want != got {
			t.Errorf("want: %s, got: %s", want, got)
		}
	})
}
