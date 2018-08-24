package exec

import (
	"time"
)

// HandleFunc handler after execution of Cmd
type HandleFunc func(output string, err error) (success bool, msg string)

// PipeCmd cmd in pipeline
type PipeCmd struct {
	Cmd        *Cmd
	HandleFunc HandleFunc
	HandleMsg  string
	DelayDur   time.Duration
	pipe       *Pipeline
}

// Pipeline HandleFunc to stop it at error
func pipeStopAtErr(s string, e error) (bool, string) {
	if e == nil {
		return true, ""
	}
	return false, e.Error()
}

// Pipeline HandleFunc to skip error
func pipeSkipErr(string, error) (bool, string) {
	return true, ""
}

// Handle set HandleFunc to cmd for handling final output or error.
func (pcmd *PipeCmd) Handle(f HandleFunc) *PipeCmd {
	pcmd.HandleFunc = f
	return pcmd
}

// Delay delay cmd execution
func (pcmd *PipeCmd) Delay(dur time.Duration) *PipeCmd {
	pcmd.DelayDur = dur
	return pcmd
}

// SkipErr add handler after cmd to skip err.
func (pcmd *PipeCmd) SkipErr(f HandleFunc) *PipeCmd {
	pcmd.HandleFunc = pipeSkipErr
	return pcmd
}

// Then define next step in pipeline
func (pcmd *PipeCmd) Then(name string, arg ...string) *PipeCmd {
	return pcmd.pipe.Do(name, arg...)
}

// Start start pipeline, pointer of pipeline.Start()
func (pcmd *PipeCmd) Start() error {
	return pcmd.pipe.Start()
}

// Run run pipeline, pointer of pipeline.Run()
func (pcmd *PipeCmd) Run() error {
	return pcmd.pipe.Run()
}
