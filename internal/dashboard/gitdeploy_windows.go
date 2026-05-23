//go:build !linux
// +build !linux

package dashboard

import (
	"os/exec"
)

func setProcAttributes(cmd *exec.Cmd) {
	// No-op on Windows/other systems
}

func killProcessGroup(cmd *exec.Cmd) {
	if cmd.Process != nil {
		_ = cmd.Process.Kill()
	}
}
