// +build !windows

package exec

import (
	"os"
	"os/exec"
	"syscall"
)

func updateCmdByOS(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
}

func killProcess(p *os.Process) {
	syscall.Kill(-p.Pid, syscall.SIGKILL)
}
