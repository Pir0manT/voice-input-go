package fastflowlm

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

// setSysProcAttr устанавливает платформенные атрибуты процесса (Windows: CREATE_NEW_PROCESS_GROUP)
func setSysProcAttr(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: 0x00000200, // CREATE_NEW_PROCESS_GROUP
	}
}

// killProcessTree убивает процесс и все его дочерние процессы (Windows: taskkill /T /F)
func killProcessTree(proc *os.Process) error {
	// taskkill /T убивает дерево процессов, /F — принудительно
	kill := exec.Command("taskkill", "/T", "/F", "/PID", fmt.Sprintf("%d", proc.Pid))
	kill.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	return kill.Run()
}
