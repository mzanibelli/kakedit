package picker

import (
	"fmt"
	"os"
	"os/exec"
)

// Picker is an external tool that can edit files.
type Picker struct {
	bin string
}

// New creates a picker.
func New(bin string) *Picker {
	return &Picker{bin}
}

// Pick runs the external tool and replaces EDITOR to allow sending edit
// instructions to the remote instance.
func (p Picker) Pick(editorCmd fmt.Stringer) *exec.Cmd {
	cmd := exec.Command(p.bin)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	for _, name := range [...]string{"EDITOR", "VISUAL"} {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", name, editorCmd))
	}

	return cmd
}
