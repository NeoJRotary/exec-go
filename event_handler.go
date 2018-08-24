package exec

// EventHandler bind event hanlder to do logging or other handling as you want
// be careful that all event handler function will be triggered synchronously to provide you synchronously control.
// use `go func()` inside handler for asynchronously control
type EventHandler struct {
	// trigger when cmd successfully started
	CmdStarted CmdEvent
	// trigger when cmd get any output ( error will trigger CmdFailed )
	CmdRead CmdEvent
	// trigger when cmd be canceled
	CmdCanceled CmdEvent
	// trigger when cmd failed which means process throw error
	CmdFailed CmdEvent
	// trigger when cmd done without any error
	CmdDone CmdEvent

	// trigger when pipeline successfully started
	PipelineStarted PipelineEvent
	// trigger when pipeline be canceled
	PipelineCanceled PipelineEvent
	// trigger when pipeline failed
	PipelineFailed PipelineEvent
	// trigger when pipeline done without any error
	PipelineDone PipelineEvent
}

// CmdEvent cmd event handle func
type CmdEvent func(cmd *Cmd)

// PipelineEvent pipeline event handle func
type PipelineEvent func(pipe *Pipeline)
