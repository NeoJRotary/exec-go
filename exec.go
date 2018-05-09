package exec

// Exec main type
type Exec struct {
	EventHandler *EventHandler
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
func NewExec(eh *EventHandler) *Exec {
	return &Exec{EventHandler: eh}
}

// Run run a cmd at current dir and wait for response
func (ex *Exec) Run(name string, arg ...string) (out string, err error) {
	cmd := NewCmd("", name, arg...)
	cmd.SetEventHandler(ex.EventHandler)
	out, err = cmd.Run()
	return out, err
}

// Dir return Run() which exec at dir
func (ex *Exec) Dir(dir string) AtDir {
	return AtDir{
		Dir: dir,
		Run: func(name string, arg ...string) (out string, err error) {
			cmd := NewCmd(dir, name, arg...)
			out, err = cmd.Run()
			cmd.SetEventHandler(ex.EventHandler)
			return out, err
		},
	}
}

// NewCmd get new Cmd which use EventHandler from Exec
func (ex *Exec) NewCmd(dir string, name string, arg ...string) *Cmd {
	return NewCmd(dir, name, arg...).SetEventHandler(ex.EventHandler)
}

// // NewPipeline get new Pipeline which use EventHandler from Exec
// func (ex *Exec) NewPipeline() *Pipeline {
// 	return NewPipeline().SetEventHandler(ex.EventHandler)
// }
