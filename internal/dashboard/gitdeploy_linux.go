//go:build linux
// +build linux

package dashboard

import (
	"os/exec"
	"syscall"
)

func setProcAttributes(cmd *exec.Cmd) {
	if cmd.SysProcAttr == nil {
		cmd.SysProcAttr = &syscall.SysProcAttr{}
	}
	cmd.SysProcAttr.Setpgid = true
	// Enforce memory limit of 2GB (Soft and Hard RSS/AS limits)
	const memLimit = 2 * 1024 * 1024 * 1024 // 2 GB
	cmd.SysProcAttr.Rlimits = []syscall.Rlimit{
		{
			Type: syscall.RLIMIT_AS,
			Cur:  memLimit,
			Max:  memLimit,
		},
		{
			Type: syscall.RLIMIT_DATA,
			Cur:  memLimit,
			Max:  memLimit,
		},
	}
}

func killProcessGroup(cmd *exec.Cmd) {
	if cmd.Process != nil {
		// Negative PID kills the process group
		_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
	}
}
