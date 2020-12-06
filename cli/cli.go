package cli

import (
	"flag"
	"fmt"
	"os"
)

var (
	// Mode is the selected operation mode.
	Mode OperationMode = ModeServer
)

func init() {
	flag.Var(&Mode, "mode",
		"select mode of operation (server, client, local)")
}

// Fail exits with an optional error and prints usage information.
func Fail(err error) {
	if err != nil {
		fmt.Fprintln(flag.CommandLine.Output(), "Error:", err)
	}
	fmt.Fprintf(flag.CommandLine.Output(),
		"Usage: %s -mode MODE PROGRAM [ARGS...]", os.Args[0])
	os.Exit(1)
}

// Parse validates CLI arguments according to the selected mode.
func Parse() error {
	flag.Parse()
	switch Mode {
	case ModeServer, ModeLocal:
		return minArgs(1)
	case ModeClient:
		return minArgs(2)
	default:
		return nil
	}
}

func minArgs(n int) error {
	if !flag.Parsed() {
		return fmt.Errorf("could not validate args before parsing them")
	}
	if flag.NArg() < n {
		return fmt.Errorf("mode %s requires at least %d arguments", Mode, n)
	}
	return nil
}

// OperationMode is a enumeration of the possible modes for this program. It is
// used as a command line flag.
type OperationMode int

const (
	// ModeServer starts listening for remote edit requests.
	ModeServer OperationMode = iota
	// ModeClient forwards edit requests to the server.
	ModeClient
	// ModeLocal makes edit requests start a new Kakoune instance.
	ModeLocal
)

var allowedModes = map[string]OperationMode{
	"server": ModeServer,
	"client": ModeClient,
	"local":  ModeLocal,
}

// Set allows the flag package to use this type as a CLI flag.
func (m *OperationMode) Set(value string) error {
	mode, ok := allowedModes[value]
	if !ok {
		return fmt.Errorf("unsupported mode: %s", value)
	}
	*m = mode
	return nil
}

// String implements the standard fmt.Stringer interface.
func (m OperationMode) String() string {
	for key, val := range allowedModes {
		if val == m {
			return key
		}
	}
	return ""
}