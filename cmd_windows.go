// +build windows

package exec

import (
	"os"
	"os/exec"
	"strconv"
)

func updateCmdByOS(cmd *exec.Cmd) {
}

func killProcess(p *os.Process) {
	exec.Command("taskkill", "/PID", strconv.Itoa(p.Pid), "/T", "/F").Run()
	p.Kill()
}
