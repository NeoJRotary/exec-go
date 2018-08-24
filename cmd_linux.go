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
	if p == nil {
		return
	}
	syscall.Kill(-p.Pid, syscall.SIGKILL)
}
