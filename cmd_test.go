package exec_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/NeoJRotary/exec-go"
)

func TestCmd_Run(t *testing.T) {
	out, err := exec.NewCmd("", "echo", "123").Run()
	if err != nil {
		t.Fatal(err)
	}
	if out != "123\n" {
		t.Fatal("wrong out, should be 123\\n, get", out)
	}

	out, err = exec.NewCmd("", "bash", "-c", "echo errrr 1>&2; exit 1;").Run()
	if err == nil {
		t.Fatal("should get error")
	}
	if err.Error() != "errrr\n" {
		t.Fatal("wrong error, should be errrr\\n, get", err.Error())
	}
}

func TestCmd_AddEnv(t *testing.T) {
	cmd := exec.NewCmd("", "bash", "-c", `echo $MYEnv`).AddEnv("MYEnv", "123")
	out, err := cmd.Run()
	if err != nil {
		t.Fatal(err)
	}
	if out != "123\n" {
		t.Fatal("wrong out, should be 123\\n, get", out)
	}
}

func TestCmd_StartWait(t *testing.T) {
	cmd := exec.NewCmd("", "echo", "123")
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
	cmd := exec.NewCmd("", "bash", "-c", `echo 123; ccccc`)
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
	cmd := exec.NewCmd("", "bash", "-c", `echo 123; sleep 1s; echo 456; printf 789`)
	err := cmd.Start()
	if err != nil {
		t.Fatal(err)
	}

	cmd.Read()
	b := string(cmd.GetMsg())
	if b != "123\n" {
		t.Fatal("wrong, get:", b)
	}
	time.Sleep(time.Millisecond * 1100)
	cmd.Read()
	b = string(cmd.GetMsg())
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
	cmd := exec.NewCmd("", "bash", "-c", `bash -c "printf 123; sleep 5s; printf 456"`)
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

func TestCmd_Timeout(t *testing.T) {
	cmd := exec.NewCmd("", "sleep", "10s").SetTimeout(time.Second * 2)
	err := cmd.Start()
	if err != nil {
		t.Fatal(err)
	}
	startAt := time.Now()
	err = cmd.Wait()
	if err == nil {
		t.Fatal("should throw err")
	}
	if !cmd.TimedOut {
		t.Fatal("should time out")
	}
	if time.Since(startAt).Seconds() > 2.5 {
		t.Fatal("process should time out in 2.5s")
	}
}
