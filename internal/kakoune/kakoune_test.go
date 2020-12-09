package kakoune_test

import (
	"kakedit/internal/kakoune"
	"os"
	"testing"
)

func TestEditCommandFormat(t *testing.T) {
	t.Skip("TODO")
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
