package exec

// DefaultEventHandler default event handler
var DefaultEventHandler = &EventHandler{}

// RunCmd run cmd and wait for response
func RunCmd(dir, name string, arg ...string) (out string, err error) {
	cmd := NewCmd(dir, name, arg...)
	return cmd.Run()
}
