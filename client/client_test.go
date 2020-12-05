package client_test

import (
	"kakedit/client"
	"testing"
)

func TestClient(t *testing.T) {
	ed := client.New("foo", "bar")
	want := "foo bar"
	got := ed.String()
	if want != got {
		t.Errorf("want: %s, got: %s", want, got)
	}
}
