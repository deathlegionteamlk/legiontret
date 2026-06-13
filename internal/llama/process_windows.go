//go:build windows

package llama

import (
	"os/exec"
)

// stopProcess stops a process on Windows
func stopProcess(cmd *exec.Cmd) error {
	if cmd == nil || cmd.Process == nil {
		return nil
	}
	return cmd.Process.Kill()
}

// setProcessGroupAttr is a no-op on Windows
func setProcessGroupAttr(cmd *exec.Cmd) {
}
