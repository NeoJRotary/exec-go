package exec

import (
	"errors"
	"strings"
	"time"
)

// Pipeline commands pipeline
type Pipeline struct {
	// pipeline started
	Started bool
	// pipeline done
	Done bool
	// pipeline failed
	Failed bool
	// pipeline canceld
	Canceled bool
	// pipline Timedout
	Timedout bool
	// Cmd pointer slice
	PipeCmds []*PipeCmd
	// current index of PipeCmd
	CurrentIndex int
	// EventHandler
	EventHandler *EventHandler
	// TimeoutDur duration
	TimeoutDur time.Duration
	// FailureMsg custom error message when failed at handler
	FailureMsg string

	// current Under
	currentDir string
	// current Env
	currentEnv []string
	// current Event Handler
	currentEventHandler *EventHandler

	// chan to check pipeline is running
	runningChan chan bool
}

// NewPipeline new pipeline.
func NewPipeline() *Pipeline {
	pipe := Pipeline{
		currentDir:   "",
		currentEnv:   []string{},
		EventHandler: DefaultEventHandler,
	}
	return &pipe
}

// SetPipelineEventHandler set event handler of pipeline
func (pipe *Pipeline) SetPipelineEventHandler(eh *EventHandler) *Pipeline {
	pipe.EventHandler = eh
	pipe.currentEventHandler = eh
	return pipe
}

// SetTimeout set timeout
func (pipe *Pipeline) SetTimeout(dur time.Duration) *Pipeline {
	pipe.TimeoutDur = dur
	return pipe
}

// Under under which dir.
func (pipe *Pipeline) Under(dir string) *Pipeline {
	if pipe.Started {
		return pipe
	}
	pipe.currentDir = dir
	return pipe
}

// SetEnv replace whole Env for cmds. KeyValue should be like 'Key=Value'.
func (pipe *Pipeline) SetEnv(keyValue string, more ...string) *Pipeline {
	if pipe.Started {
		return pipe
	}
	pipe.currentEnv = append(more, keyValue)
	return pipe
}

// SetCmdEventHandler set Event Handler for cmds. It inheirt from pipeline's by default.
func (pipe *Pipeline) SetCmdEventHandler(eh *EventHandler) *Pipeline {
	pipe.currentEventHandler = eh
	return pipe
}

// AddEnv Add Env.
func (pipe *Pipeline) AddEnv(name, value string) *Pipeline {
	if pipe.Started {
		return pipe
	}
	pipe.currentEnv = append(pipe.currentEnv, name+"="+value)
	return pipe
}

// RemoveEnv Remove Env by name.
func (pipe *Pipeline) RemoveEnv(name string) *Pipeline {
	if pipe.Started {
		return pipe
	}
	newEnv := []string{}
	for _, Env := range pipe.currentEnv {
		if !strings.HasPrefix(Env, name+"=") {
			newEnv = append(newEnv, Env)
		}
	}
	pipe.currentEnv = newEnv
	return pipe
}

// Do add cmd to pipeline. It use latest Dir/Env you set before.
func (pipe *Pipeline) Do(name string, arg ...string) *PipeCmd {
	if pipe.Started {
		return nil
	}

	cmd := NewCmd(pipe.currentDir, name, arg...).SetEnv(pipe.currentEnv)
	if pipe.currentEventHandler != nil {
		cmd.SetEventHandler(pipe.currentEventHandler)
	} else if pipe.EventHandler != nil {
		cmd.SetEventHandler(pipe.EventHandler)
	}

	pcmd := &PipeCmd{
		Cmd:        cmd,
		HandleFunc: pipeStopAtErr,
		pipe:       pipe,
	}
	pipe.PipeCmds = append(pipe.PipeCmds, pcmd)
	return pcmd
}

// Finished pipeline lifecycle is finished or not
func (pipe *Pipeline) Finished() bool {
	return pipe.Done || pipe.Failed || pipe.Canceled
}

// Start start pipeline
func (pipe *Pipeline) Start() error {
	if pipe.Started {
		return errors.New("Pipeline already started")
	}
	pipe.Started = true
	pipe.runningChan = make(chan bool, 1)
	go pipe.runCmds()
	if pipe.TimeoutDur != 0 {
		go pipe.checkTimeout()
	}
	if pipe.EventHandler.PipelineStarted != nil {
		pipe.EventHandler.PipelineStarted(pipe)
	}
	return nil
}

// do timeout checking
func (pipe *Pipeline) checkTimeout() {
	timer := time.NewTimer(pipe.TimeoutDur)
	<-timer.C
	if !(pipe.Failed || pipe.Done) && !pipe.Canceled {
		pipe.Timedout = true
		pipe.Cancel()
	}
}

func (pipe *Pipeline) runCmds() {
	for i := 0; i < len(pipe.PipeCmds); i++ {
		if pipe.Finished() {
			break
		}

		c := pipe.PipeCmds[i]

		time.Sleep(c.DelayDur)
		if pipe.Finished() {
			break
		}

		pipe.CurrentIndex = i
		err := c.Cmd.Start()
		if err == nil {
			err = c.Cmd.Wait()
		}

		if pipe.Canceled {
			break
		} else {
			success, msg := c.HandleFunc(c.Cmd.Output(), err)
			c.HandleMsg = msg
			if !success {
				pipe.Failed = true
				pipe.FailureMsg = msg
				break
			}
		}
	}
	if !pipe.Canceled && !pipe.Failed {
		pipe.Done = true
	}

	pipe.runningChan <- false
}

// Wait block until Done/Failed/Canceled
func (pipe *Pipeline) Wait() error {
	if !pipe.Started {
		return errors.New("Pipeline cant Wait() before Start()")
	}
	if pipe.runningChan == nil {
		return errors.New("Pipeline cant Wait() after lifecycle is finished")
	}

	<-pipe.runningChan
	close(pipe.runningChan)
	pipe.runningChan = nil

	if pipe.Failed && pipe.EventHandler.PipelineFailed != nil {
		pipe.EventHandler.PipelineFailed(pipe)
	}
	if pipe.Done && pipe.EventHandler.PipelineDone != nil {
		pipe.EventHandler.PipelineDone(pipe)
	}

	return nil
}

// Run run pipeline and wait for result
func (pipe *Pipeline) Run() error {
	err := pipe.Start()
	if err != nil {
		return err
	}
	return pipe.Wait()
}

// Cancel cancel pipeline
func (pipe *Pipeline) Cancel() {
	if pipe.Canceled {
		return
	}
	pipe.Canceled = true
	pipe.PipeCmds[pipe.CurrentIndex].Cmd.Cancel()
	if pipe.EventHandler.PipelineCanceled != nil {
		pipe.EventHandler.PipelineCanceled(pipe)
	}
}

// GetCmdsOutput get output of all Cmds concat into one string
func (pipe *Pipeline) GetCmdsOutput() (string, error) {
	if !pipe.Finished() {
		return "", errors.New("Pipeline cant GetCmdsOutput() before lifecycle is finished")
	}

	strs := []string{}
	for i := 0; i <= pipe.CurrentIndex; i++ {
		c := pipe.PipeCmds[i].Cmd
		strs = append(strs, c.Output()+c.Error())
	}
	return strings.Join(strs, ""), nil
}
