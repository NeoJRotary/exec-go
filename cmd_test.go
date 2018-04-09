package exec

import (
	"fmt"
	"testing"
	"time"
)

func TestCmd_Run(t *testing.T) {
	cmd := NewCmd("", "echo", "123")
	out, err := cmd.Run()
	if err != nil {
		t.Fatal(err)
	}
	if out != "123\n" {
		t.Fatal("wrong out, should be 123\\n, get", out)
	}
}

func TestCmd_AddEnv(t *testing.T) {
	cmd := NewCmd("", "bash", "-c", `echo $MYENV`).AddEnv("MYENV", "123")
	out, err := cmd.Run()
	if err != nil {
		t.Fatal(err)
	}
	if out != "123\n" {
		t.Fatal("wrong out, should be 123\\n, get", out)
	}
}

func TestCmd_StartWait(t *testing.T) {
	cmd := NewCmd("", "echo", "123")
	err := cmd.Start()
	if err != nil {
		t.Fatal(err)
	}
	err = cmd.Wait()
	if err != nil {
		t.Fatal(err, cmd.Error())
	}
	if cmd.Output() != "123\n" {
		t.Fatal("wrong out, should be 123\\n")
	}
}

func TestCmd_Output(t *testing.T) {
	cmd := NewCmd("", "bash", "-c", `echo 123; ccccc`)
	err := cmd.Start()
	if err != nil {
		t.Fatal(err)
	}
	err = cmd.Wait()
	if err == nil {
		t.Fatal("should throw error")
	}

	if cmd.Output() != "123\n" {
		t.Fatal("wrong output, get:", cmd.Output())
	}
	if cmd.Error() != "bash: ccccc: command not found\n" {
		fmt.Println([]byte(cmd.Error()))
		fmt.Println([]byte("bash: ccccc: command not found"))
		t.Fatal("wrong output, get:", cmd.Error())
	}
}

func TestCmd_Read(t *testing.T) {
	cmd := NewCmd("", "bash", "-c", `echo 123; sleep 1s; echo 456; printf 789`)
	err := cmd.Start()
	if err != nil {
		t.Fatal(err)
	}

	cmd.Read()
	b := string(cmd.Get())
	if b != "123\n" {
		t.Fatal("wrong, get:", b)
	}
	time.Sleep(time.Millisecond * 1100)
	cmd.Read()
	b = string(cmd.Get())
	if b != "456\n789" {
		t.Fatal("wrong, get:", b)
	}

	err = cmd.Wait()
	if err != nil {
		t.Fatal(err)
	}
	if cmd.Output() != "123\n456\n789" {
		t.Error("wrong output, get:\n", cmd.Output())
	}
}

func TestCmd_Cancel(t *testing.T) {
	cmd := NewCmd("", "bash", "-c", `bash -c "printf 123; sleep 5s; printf 456"`)
	err := cmd.Start()
	if err != nil {
		t.Fatal(err)
	}
	startAt := time.Now()
	go func() {
		time.Sleep(time.Millisecond * 500)
		cmd.Cancel()
	}()
	err = cmd.Wait()
	if err == nil {
		t.Fatal("should throw err")
	}
	if !cmd.Canceled {
		t.Fatal("Canceled should be true")
	}
	if time.Since(startAt).Seconds() >= 1 {
		t.Fatal("process should stop before 1s")
	}
	if cmd.Output() != "123" {
		t.Fatal("output should be 123")
	}
}
