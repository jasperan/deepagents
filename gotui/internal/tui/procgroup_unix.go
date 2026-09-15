//go:build unix

package tui

import (
	"os/exec"
	"syscall"
)

// setProcessGroup puts the child in its own process group so a cancellation can
// take down whatever it spawned.
//
// The deepagents CLI starts a `langgraph dev` server subprocess. Killing only
// the direct child would leave that grandchild holding the stdout pipe open, so
// the reader goroutine would block until it exited on its own and the UI would
// look frozen long after the timeout.
func setProcessGroup(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
}

// killProcessGroup signals the whole group, which closes inherited pipes
// immediately instead of waiting for the grandchildren to exit.
func killProcessGroup(cmd *exec.Cmd) error {
	if cmd.Process == nil {
		return nil
	}
	if err := syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL); err != nil {
		return cmd.Process.Kill()
	}
	return nil
}
