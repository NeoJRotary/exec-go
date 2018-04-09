package exec

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
	"os/exec"
	"strings"
)

// Cmd ...
type Cmd struct {
	Dir       string
	Name      string
	Args      []string
	Started   bool
	Canceled  bool
	Failed    bool
	Done      bool
	cmd       *exec.Cmd
	outPipe   io.ReadCloser
	outReader *bufio.Reader
	errBuf    *bytes.Buffer
	out       []byte
	err       []byte
	line      []byte
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
		Dir:    dir,
		Name:   name,
		Args:   arg,
		cmd:    cmd,
		errBuf: &errBuf,
	}
}

// AddEnv add env
func (c *Cmd) AddEnv(name, value string) *Cmd {
	c.cmd.Env = append(c.cmd.Env, name+"="+value)
	return c
}

// GetCmd get cmd string
func (c *Cmd) GetCmd() string {
	return c.Name + " " + strings.Join(c.Args, " ")
}

// Run run cmd. If error, it return output until error and error messages as error.
func (c *Cmd) Run() (string, error) {
	var outBuf bytes.Buffer
	c.cmd.Stdout = &outBuf
	err := c.cmd.Run()
	c.out = outBuf.Bytes()
	c.err = c.errBuf.Bytes()
	if err != nil {
		return outBuf.String(), errors.New(c.errBuf.String())
	}
	return outBuf.String(), nil
}

// Start start cmd
func (c *Cmd) Start() error {
	outPipe, _ := c.cmd.StdoutPipe()
	c.outPipe = outPipe

	err := c.cmd.Start()
	if err != nil {
		return err
	}
	c.Started = true
	c.outReader = bufio.NewReader(outPipe)
	return nil
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
	c.line = p[:n]
	return true
}

// Get get last output from cmd
func (c *Cmd) Get() []byte {
	return c.line
}

// Cancel terminate the command process
func (c *Cmd) Cancel() {
	c.Canceled = true
	c.Read()
	killProcess(c.cmd.Process)
}

// Wait wait for cmd until its done. It will return original error ( exit(1) ) except error message. Use Error() to get error message.
func (c *Cmd) Wait() error {
	for c.Read() {
	}
	err := c.cmd.Wait()
	if err != nil {
		c.Failed = true
	} else {
		c.Done = true
	}
	c.err = c.errBuf.Bytes()
	return err
}

// Output get outputs of cmd (without error) in bytes
func (c *Cmd) Output() string {
	return string(c.out)
}

// Error get error output in bytes. Return nil if no error.
func (c *Cmd) Error() string {
	return string(c.err)
}
