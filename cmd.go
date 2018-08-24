package exec

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"
)

// Cmd ...
type Cmd struct {
	// run at dir
	Dir string
	// target executable name
	Name string
	// arguments of cmd
	Args []string
	// exec started
	Started bool
	// exec canceled
	Canceled bool
	// exec timedout
	TimedOut bool
	// exec failed by error
	Failed bool
	// exec done
	Done bool
	// exec timeout
	Timeout time.Duration
	// exec duration
	Duration time.Duration
	// EventHandler pointer
	EventHandler *EventHandler
	// native exec.Cmd target
	Cmd *exec.Cmd
	// start time
	startTime time.Time
	// output pipe
	outPipe io.ReadCloser
	// output buf reader
	outReader *bufio.Reader
	// err buf
	errBuf *bytes.Buffer
	// result output bytes
	out []byte
	// result error bytes
	err []byte
	// latest ouput
	msg []byte
}

// NewCmd new Cmd object
func NewCmd(dir string, name string, arg ...string) *Cmd {
	cmd := exec.Command(name, arg...)
	if dir == "" {
		dir = "./"
	}
	cmd.Dir = dir
	cmd.Env = os.Environ()
	updateCmdByOS(cmd)
	var errBuf bytes.Buffer
	cmd.Stderr = &errBuf
	return &Cmd{
		Dir:          dir,
		Name:         name,
		Args:         arg,
		Cmd:          cmd,
		EventHandler: DefaultEventHandler,
		errBuf:       &errBuf,
	}
}

// SetEventHandler set event handler
func (c *Cmd) SetEventHandler(eh *EventHandler) *Cmd {
	c.EventHandler = eh
	return c
}

// SetEnv set Envs
func (c *Cmd) SetEnv(Env []string) *Cmd {
	c.Cmd.Env = Env
	return c
}

// AddEnv add Env
func (c *Cmd) AddEnv(name, value string) *Cmd {
	c.Cmd.Env = append(c.Cmd.Env, name+"="+value)
	return c
}

// SetTimeout set timeout
func (c *Cmd) SetTimeout(dur time.Duration) *Cmd {
	c.Timeout = dur
	return c
}

// GetCmd get cmd string
func (c *Cmd) GetCmd() string {
	return c.Name + " " + strings.Join(c.Args, " ")
}

// Start start cmd
func (c *Cmd) Start() error {
	outPipe, _ := c.Cmd.StdoutPipe()
	c.outPipe = outPipe

	err := c.Cmd.Start()
	if err != nil {
		return err
	}
	if c.Timeout != 0 {
		go c.checkTimeout()
	}
	if c.EventHandler.CmdStarted != nil {
		c.EventHandler.CmdStarted(c)
	}
	c.startTime = time.Now()
	c.Started = true
	c.outReader = bufio.NewReader(outPipe)
	return nil
}

// do timeout checking
func (c *Cmd) checkTimeout() {
	timer := time.NewTimer(c.Timeout)
	<-timer.C
	if !(c.Failed || c.Done) && !c.Canceled {
		c.TimedOut = true
		c.Cancel()
	}
}

// Read read latest output from last Read()
func (c *Cmd) Read() bool {
	if c.Canceled || c.Done {
		return false
	}
	p := make([]byte, 65536)
	n, err := c.outReader.Read(p)
	if err != nil {
		return false
	}
	c.out = append(c.out, p[:n]...)
	c.msg = p[:n]
	if c.EventHandler.CmdRead != nil {
		c.EventHandler.CmdRead(c)
	}
	return true
}

// GetMsg get last output from cmd
func (c *Cmd) GetMsg() []byte {
	return c.msg
}

// Cancel terminate the command process
func (c *Cmd) Cancel() {
	if c.Canceled {
		return
	}
	c.Canceled = true
	c.Read()
	killProcess(c.Cmd.Process)
	if c.EventHandler.CmdCanceled != nil {
		c.EventHandler.CmdCanceled(c)
	}
}

// Wait wait for cmd until its done. It will return original error ( like exit(1) ) except error message. Use Error() to get error message.
func (c *Cmd) Wait() error {
	for c.Read() {
	}

	err := c.Cmd.Wait()
	if !c.Canceled {
		if err != nil {
			c.Failed = true
		} else {
			c.Done = true
		}
	}

	c.Duration = time.Since(c.startTime)
	c.err = c.errBuf.Bytes()

	if c.Failed && c.EventHandler.CmdFailed != nil {
		c.EventHandler.CmdFailed(c)
	}
	if c.Done && c.EventHandler.CmdDone != nil {
		c.EventHandler.CmdDone(c)
	}
	return err
}

// Run run cmd. If error, it return output until error and error messages as error.
func (c *Cmd) Run() (string, error) {
	err := c.Start()
	if err != nil {
		return "", err
	}
	err = c.Wait()
	if err != nil {
		return c.Output(), errors.New(c.Error())
	}
	return c.Output(), nil
}

// Output get outputs of cmd (without error) in string
func (c *Cmd) Output() string {
	return string(c.out)
}

// Error get error output in string
func (c *Cmd) Error() string {
	return string(c.err)
}
