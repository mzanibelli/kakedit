package cli_test

import (
	"fmt"
	"kakedit/cli"
	"os"
	"strings"
	"testing"
)

func TestAllowedModes(t *testing.T) {
	tests := [...]struct {
		val string
		ok  bool
	}{
		{"foo", false},
		{"client", true},
		{"server", true},
		{"local", true},
	}
	for _, test := range tests {
		t.Run(test.val, func(t *testing.T) {
			mode := new(cli.OperationMode)
			err := mode.Set(test.val)
			if (err == nil) != test.ok {
				t.Log(err)
				t.Errorf("should work: %t", test.ok)
			}
		})
	}
}

func TestModePrint(t *testing.T) {
	tests := [...]struct {
		mode cli.OperationMode
		val  string
	}{
		{cli.ModeServer, "server"},
		{cli.ModeClient, "client"},
		{cli.ModeLocal, "local"},
	}
	for _, test := range tests {
		t.Run(test.val, func(t *testing.T) {
			got := fmt.Sprintf("%s", test.mode)
			if got != test.val {
				t.Errorf("want: %s, got: %s", test.val, got)
			}
		})
	}
}

func TestParse(t *testing.T) {
	tests := [...]struct {
		name string
		args string
		ok   bool
	}{
		{"server 1 args", "kakedit -mode server /bin/true", true},
		{"server 2 args", "kakedit -mode server /bin/true hello", true},
		{"server 0 args", "kakedit -mode server", false},
		{"local 1 args", "kakedit -mode local /bin/true", true},
		{"local 2 args", "kakedit -mode local /bin/true hello", true},
		{"local 0 args", "kakedit -mode local", false},
		{"client 1 args", "kakedit -mode client /bin/true", false},
		{"client 2 args", "kakedit -mode client /bin/true hello", true},
		{"client 0 args", "kakedit -mode client", false},
	}
	for _, test := range tests {
		t.Run(test.args, func(t *testing.T) {
			os.Args = strings.Fields(test.args)
			err := cli.Parse()
			if (err == nil) != test.ok {
				t.Log(err)
				t.Errorf("should work: %t", test.ok)
			}
		})
	}
}
