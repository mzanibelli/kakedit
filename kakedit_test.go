package kakedit_test

import (
	"errors"
	"fmt"
	"kakedit"
	"os"
	"os/exec"
	"path"
	"syscall"
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

	pwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	session := fmt.Sprint(os.Getpid())
	client := "client0"

	kak := exec.Command(bin, "-ui", "dummy", "-n", "-s", session, "-e",
		fmt.Sprintf("'rename-client %s'", client))

	kak.Env = os.Environ()

	if err := kak.Start(); err != nil {
		t.Fatal(err)
	}

	defer func() {
		kak.Process.Signal(syscall.SIGTERM)
		kak.Wait()
	}()

	pick := path.Join(pwd, "testdata", "pick.sh")
	pipe := path.Join(pwd, "testdata", "pipe.sh")

	tests := [...]struct {
		name    string
		pick    string
		pipe    string
		wrap    string
		session string
		client  string
		err     error
	}{
		{"server simple", "/bin/true", "/bin/true", "/bin/true", session, client, nil},
		{"server cmd fail", "/bin/false", "/bin/true", "/bin/true", session, client, errors.New("exit status 1")},
		{"server pipe fail", pick, "/bin/false", "/bin/true", session, client, errors.New("exit status 1")},
		{"server roundtrip", pick, pipe, "/bin/true", session, client, nil},
		{"server kak fail", pick, pipe, "/bin/true", "unknown", "unknown", errors.New("exit status 255")},
		{"local simple", "/bin/true", "/bin/true", "/bin/true", "", "", nil},
		{"local wrap fail", pick, "/bin/true", "/bin/false", "", "", errors.New("exit status 1")},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := os.Setenv("kak_session", test.session); err != nil {
				t.Fatal(err)
			}
			if err := os.Setenv("kak_client", test.client); err != nil {
				t.Fatal(err)
			}

			err := kakedit.ExternalProgram(test.pick, test.pipe, test.wrap)

			want := fmt.Sprint(test.err)
			got := fmt.Sprint(err)
			if want != got {
				t.Errorf("want: %s, got: %s", want, got)
			}
		})
	}
}
