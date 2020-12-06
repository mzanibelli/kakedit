package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"
	"testing"
)

// /!\ Experimental integration test.
func TestKakEdit(t *testing.T) {
	// We are going to run a headless kakoune session for this test...
	// Ensure we are in a good mood for this beast.
	bin, err := exec.LookPath("kak")
	if err != nil {
		t.Skip("missing kak(1) executable")
	}
	if testing.Short() {
		t.Skip("skipping integration test when short flag is given")
	}

	// This program invokes itself, so the integration test must
	// use a compiled binary. To prevent any error related to the
	// use of a compiled program older than the current code, force
	// recompiling before running the test.
	pwd, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}
	self := path.Join(pwd, "..", "kakedit")
	build := exec.Command("go", "build", "-o", self, path.Join(pwd, "kakedit.go"))
	if out, err := build.CombinedOutput(); err != nil {
		t.Logf("%s", out)
		t.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	session := fmt.Sprint(os.Getpid())
	client := "client0"

	kak := exec.CommandContext(ctx, bin, "-ui", "dummy", "-n",
		"-s", session, "-e", fmt.Sprintf("'rename-client %s'", client))
	kak.Env = os.Environ()
	go kak.Run()

	cmd := exec.Command(self, path.Join(pwd, "..", "testdata", "picker.sh"))
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, fmt.Sprintf("kak_session=%s", session))
	cmd.Env = append(cmd.Env, fmt.Sprintf("kak_client=%s", client))
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Errorf("%s", out)
	}
}
