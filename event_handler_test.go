package exec_test

import (
	"testing"
	"time"

	exec "github.com/NeoJRotary/exec-go"
)

func TestEventHandler(t *testing.T) {
	temp := ""
	eh := &exec.EventHandler{
		CmdStarted: func(cmd *exec.Cmd) {
			temp += "CmdStarted"
		},
		CmdRead: func(cmd *exec.Cmd) {
			temp += "CmdRead"
		},
		CmdCanceled: func(cmd *exec.Cmd) {
			temp += "CmdCanceled"
		},
		CmdFailed: func(cmd *exec.Cmd) {
			temp += "CmdFailed"
		},
		CmdDone: func(cmd *exec.Cmd) {
			temp += "CmdDone"
		},
		PipelineStarted: func(cmd *exec.Pipeline) {
			temp += "PipelineStarted"
		},
		PipelineCanceled: func(cmd *exec.Pipeline) {
			temp += "PipelineCanceled"
		},
		PipelineFailed: func(cmd *exec.Pipeline) {
			temp += "PipelineFailed"
		},
		PipelineDone: func(cmd *exec.Pipeline) {
			temp += "PipelineDone"
		},
	}

	exec.NewCmd("", "echo", "123").SetEventHandler(eh).Run()
	if temp != "CmdStarted"+"CmdRead"+"CmdDone" {
		t.Fatal("wrong: ", temp)
	}

	temp = ""
	exec.NewCmd("", "echoddd", "123").SetEventHandler(eh).Run()
	if temp != "" {
		t.Fatal("wrong: ", temp)
	}

	temp = ""
	exec.NewCmd("", "bash", "wedwedwwe").SetEventHandler(eh).Run()
	if temp != "CmdStarted"+"CmdFailed" {
		t.Fatal("wrong: ", temp)
	}

	temp = ""
	exec.NewCmd("", "bash", "-c", "echo 1; sleep 0.2; echo 2;").SetEventHandler(eh).Run()
	if temp != "CmdStarted"+"CmdRead"+"CmdRead"+"CmdDone" {
		t.Fatal("wrong: ", temp)
	}

	temp = ""
	cmd := exec.NewCmd("", "bash", "-c", "echo 1; sleep 1; echo 2;").SetEventHandler(eh)
	cmd.Start()
	go func() {
		time.Sleep(time.Millisecond * 500)
		cmd.Cancel()
	}()
	cmd.Wait()
	if temp != "CmdStarted"+"CmdRead"+"CmdCanceled" {
		t.Fatal("wrong: ", temp)
	}

	temp = ""
	exec.NewPipeline().
		SetPipelineEventHandler(eh).
		Do("echo", "1").
		Then("echo", "1").
		Run()
	if temp != "PipelineStarted"+"CmdStarted"+"CmdRead"+"CmdDone"+"CmdStarted"+"CmdRead"+"CmdDone"+"PipelineDone" {
		t.Fatal("wrong: ", temp)
	}

	temp = ""
	exec.NewPipeline().
		SetPipelineEventHandler(eh).
		Do("echo", "1").
		Then("bash", "qwdwqdqwdwq").
		Run()
	if temp != "PipelineStarted"+"CmdStarted"+"CmdRead"+"CmdDone"+"CmdStarted"+"CmdFailed"+"PipelineFailed" {
		t.Fatal("wrong: ", temp)
	}

	temp = ""
	pipe := exec.NewPipeline().SetPipelineEventHandler(eh)
	pipe.Do("echo", "1").
		Then("sleep", "10").
		Start()
	go func() {
		time.Sleep(time.Second)
		pipe.Cancel()
	}()
	pipe.Wait()
	if temp != "PipelineStarted"+"CmdStarted"+"CmdRead"+"CmdDone"+"CmdStarted"+"CmdCanceled"+"PipelineCanceled" {
		t.Fatal("wrong: ", temp)
	}
}
