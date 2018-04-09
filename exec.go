package exec

import (
	"log"
)

// Exec main type
type Exec struct {
	logger *log.Logger
}

// Result executed result
type Result struct {
	cmd    string
	output string
}

// AtDir contain Run() at Dir
type AtDir struct {
	Dir string
	Run func(string, ...string) (string, error)
}

// NewExec get new Exec. Give it an log.Logger for logging.
func NewExec(logger *log.Logger) *Exec {
	return &Exec{logger: logger}
}

// Run run a cmd at current dir and wait for response
func (ex *Exec) Run(name string, arg ...string) (out string, err error) {
	cmd := NewCmd("", name, arg...)
	out, err = cmd.Run()
	if ex.logger != nil {
		if err != nil {
			ex.logger.Println(cmd.GetCmd(), "\n", out, err.Error())
		} else {
			ex.logger.Println(cmd.GetCmd(), "\n", out)
		}
	}
	return out, err
}

// Dir return Run() which exec at dir
func (ex *Exec) Dir(dir string) AtDir {
	return AtDir{
		Dir: dir,
		Run: func(name string, arg ...string) (out string, err error) {
			cmd := NewCmd(dir, name, arg...)
			out, err = cmd.Run()
			if ex.logger != nil {
				if err != nil {
					ex.logger.Println(dir, cmd.GetCmd(), "\n", out, err.Error())
				} else {
					ex.logger.Println(dir, cmd.GetCmd(), "\n", out)
				}
			}
			return out, err
		},
	}
}
