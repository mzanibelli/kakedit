package kakedit_test

import (
	"context"
	"errors"
	"fmt"
	"kakedit"
	"os"
	"os/exec"
	"path"
	"testing"
)

func TestKakEdit(t *testing.T) {
	bin, err := exec.LookPath("kak")
	if err != nil {
		t.Skip("missing kak(1) executable")
	}

	// pipe.sh is a pure shell implementation of kakpipe that
	// needs netcat.
	_, err = exec.LookPath("nc")
	if err != nil {
		t.Skip("missing nc(1) executable")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	session := fmt.Sprint(os.Getpid())
	client := "client0"

	kak := exec.CommandContext(ctx, bin, "-ui", "dummy", "-n",
		"-s", session, "-e", fmt.Sprintf("'rename-client %s'", client))
	kak.Env = os.Environ()
	go kak.Run()

	pwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	pick := path.Join(pwd, "testdata", "pick.sh")
	pipe := path.Join(pwd, "testdata", "pipe.sh")

	tests := [...]struct {
		name    string
		pick    string
		pipe    string
		session string
		client  string
		err     error
	}{
		{"server simple", "/bin/true", "/bin/true", session, client, nil},
		{"server cmd fail", "/bin/false", "/bin/true", session, client, errors.New("exit status 1")},
		{"server pipe fail", pick, "/bin/false", session, client, errors.New("exit status 1")},
		{"server roundtrip", pick, pipe, session, client, nil},
		{"server kak fail", pick, pipe, "unknown", "unknown", errors.New("exit status 255")},
		{"local simple", "/bin/true", "/bin/true", "", "", nil},
		{"local cmd fail", "/bin/false", "/bin/true", "", "", errors.New("exit status 1")},
		{"local pipe fail", pick, "/bin/false", "", "", errors.New("exit status 255")},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := os.Setenv("kak_session", test.session); err != nil {
				t.Fatal(err)
			}
			if err := os.Setenv("kak_client", test.client); err != nil {
				t.Fatal(err)
			}

			err := kakedit.Run(test.pick, test.pipe)

			want := fmt.Sprint(test.err)
			got := fmt.Sprint(err)
			if want != got {
				t.Errorf("want: %s, got: %s", want, got)
			}
		})
	}
}
