package editor_test

import (
	"kakedit/editor"
	"testing"
)

func TestEditor(t *testing.T) {
	ed := editor.New("foo", "bar")
	want := "foo bar"
	got := ed.String()
	if want != got {
		t.Errorf("want: %s, got: %s", want, got)
	}
}
