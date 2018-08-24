package exec_test

import (
	"testing"

	exec "github.com/NeoJRotary/exec-go"
)

func TestRunCmd(t *testing.T) {
	r, err := exec.RunCmd("", "echo", "123")
	if err != nil {
		t.Fatal(err)
	}
	if r != "123\n" {
		t.Fatal("wrong output :", r)
	}
}
