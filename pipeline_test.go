package exec_test

import (
	"testing"
	"time"

	"github.com/NeoJRotary/exec-go"
)

func TestPipe_Basic(t *testing.T) {
	defer exec.NewCmd("", "rm", "-rf", "testDir").Run()

	pipe := exec.NewPipeline()
	pipe.SetEnv("AAA=1", "BBB=2")
	pipe.Do("mkdir", "testDir/path").Handle(func(out string, err error) (bool, string) {
		if err != nil {
			return true, "know error let it pass"
		}
		return false, "mkdir should throw error"
	})
	pipe.Do("mkdir", "-p", "testDir/path")
	pipe.Under("./testDir/path")
	pipe.Do("bash", "-c", `echo "text$AAA$BBB$CCC" > file.txt`).Delay(time.Second)
	pipe.Do("cat", "file.txt").Handle(func(out string, err error) (bool, string) {
		if err != nil {
			return false, err.Error()
		}
		if out != "text12\n" {
			return false, "should be text12, get: " + out
		}
		return true, ""
	})
	pipe.AddEnv("CCC", "3")
	pipe.RemoveEnv("BBB")
	pipe.Do("bash", "-c", "echo $AAA$BBB$CCC >> file.txt")
	pipe.Do("cat", "file.txt").Handle(func(out string, err error) (bool, string) {
		if err != nil {
			return false, err.Error()
		}
		if out != "text12\n13\n" {
			return false, "should be text1213, get: " + out
		}
		return true, ""
	})

	err := pipe.Run()
	if err != nil {
		t.Fatal(err)
	}
	if pipe.Failed {
		t.Fatal(pipe.FailureMsg)
	}

	err = pipe.Start()
	if err == nil {
		t.Fatal("should not allow to Start")
	}
	err = pipe.Wait()
	if err == nil {
		t.Fatal("should not allow to Wait")
	}

	// just run it to see there is panic or not
	pipe.GetCmdsOutput()
}

func TestPipe_Cancel(t *testing.T) {
	count := 0
	hand := func(out string, err error) (bool, string) {
		count++
		return true, ""
	}

	pipe := exec.NewPipeline()
	pipe.Do("echo", "123").Delay(time.Second).Handle(hand)
	pipe.Do("echo", "123").Delay(time.Second).Handle(hand)
	pipe.Do("sleep", "1").Handle(hand)
	pipe.Do("sleep", "1").Handle(hand)
	pipe.Do("sleep", "1").Handle(hand)

	startAt := time.Now()
	err := pipe.Start()
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		time.Sleep(time.Millisecond * 2500)
		pipe.Cancel()
	}()

	err = pipe.Wait()
	if err != nil {
		t.Fatal(err)
	}
	if !pipe.Canceled {
		t.Fatal("should be canceled")
	}
	if time.Since(startAt).Seconds() > 3 {
		t.Fatal("should cancel in 3s")
	}
	if count != 2 {
		t.Fatal("should count 2 times")
	}
}

func TestPipe_Timeout(t *testing.T) {
	count := 0
	hand := func(out string, err error) (bool, string) {
		count++
		return true, ""
	}

	pipe := exec.NewPipeline().SetTimeout(time.Millisecond * 2500)
	pipe.Do("echo", "123").Delay(time.Second).Handle(hand)
	pipe.Do("echo", "123").Delay(time.Second).Handle(hand)
	pipe.Do("sleep", "1").Handle(hand)
	pipe.Do("sleep", "1").Handle(hand)

	startAt := time.Now()
	err := pipe.Start()
	if err != nil {
		t.Fatal(err)
	}
	err = pipe.Wait()
	if err != nil {
		t.Fatal(err)
	}
	if !pipe.Timedout {
		t.Fatal("should be timeout")
	}
	if !pipe.Canceled {
		t.Fatal("should be canceled")
	}
	if time.Since(startAt).Seconds() > 3 {
		t.Fatal("should cancel in 3s")
	}
	if count != 2 {
		t.Fatal("should count 2 times")
	}
}
