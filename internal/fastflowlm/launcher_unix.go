//go:build !windows

package fastflowlm

import (
	"os"
	"os/exec"
	"syscall"
)

// setSysProcAttr устанавливает платформенные атрибуты процесса (Unix: Setpgid)
func setSysProcAttr(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
}

// killProcessTree убивает процесс и всю группу (Unix: kill -PGID)
func killProcessTree(proc *os.Process) error {
	// Убиваем всю группу процессов
	return syscall.Kill(-proc.Pid, syscall.SIGKILL)
}
