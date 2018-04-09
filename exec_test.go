package exec

import (
	"log"
	"os"
	"testing"
)

var testLogger = log.New(os.Stdout, "\n[EXEC] ", log.LstdFlags)

func TestExec_Run(t *testing.T) {
	exec := NewExec(testLogger)
	r, _ := exec.Run("echo", "123")
	if r != "123\n" {
		t.Fatal("wrong output :", r)
	}
}

func TestExec_Dir(t *testing.T) {
	exec := NewExec(testLogger)
	r, _ := exec.Dir("./").Run("ls", "Dockerfile")
	if r != "Dockerfile\n" {
		t.Fatal("wrong output :", r)
	}
}
