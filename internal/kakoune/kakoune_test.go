package kakoune_test

import (
	"kakedit/internal/kakoune"
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

	var want, got string

	want = strings.Join([]string{
		"echo",
		"'evaluate-commands -verbatim -client bar edit \"test.txt\"'",
		"| kak -p foo",
	}, " ")
	got = kakoune.FromEnvironment().EditClient("test.txt")
	if want != got {
		t.Errorf("want: %s, got: %s", want, got)
	}

	want = "kak -c foo"
	got = kakoune.FromEnvironment().EditSession()
	if want != got {
		t.Errorf("want: %s, got: %s", want, got)
	}

	if err := os.Setenv("kak_session", ""); err != nil {
		t.Error(err)
	}
	if err := os.Setenv("kak_client", ""); err != nil {
		t.Error(err)
	}

	want = "kak"
	got = kakoune.FromEnvironment().EditSession()
	if want != got {
		t.Errorf("want: %s, got: %s", want, got)
	}
}

func TestUnknownRemote(t *testing.T) {
	tests := []struct {
		session string
		client  string
		want    bool
	}{
		{"foo", "bar", false},
		{"", "bar", true},
		{"foo", "", true},
		{"", "", true},
	}
	for _, test := range tests {
		t.Run(test.session+"-"+test.client, func(t *testing.T) {
			if err := os.Setenv("kak_session", test.session); err != nil {
				t.Error(err)
			}
			if err := os.Setenv("kak_client", test.client); err != nil {
				t.Error(err)
			}
			got := kakoune.FromEnvironment().UnknownRemote()
			if got != test.want {
				t.Errorf("want: %t, got: %t", test.want, got)
			}
		})
	}
}
