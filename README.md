# Exec-Go v0.2.1-alpha
External command exec library for golang   

## Check details at [GoDoc](https://godoc.org/github.com/NeoJRotary/describe-go) 

## Features
- Util of native `os.exec.cmd` library.
- Pipeline library for commands.
- (Planning) PipelineGroup for parallelly manage multi-pipelines.
- (Plan to Deprecate) Cron Job function for command.

## Cmd
- Error Message   
  Error from native library is the output of `exit`, you will get error like `exit status 1` except real error message. Use `cmd.Error()` to get it.
- Cancel   
  For both Linux and Windows. On Linux, it will kill process and it's children. On Windows, it may not work if the process create child processes.
- STDOUT Streaming   
  You can do it by `Read()` and `GetMsg()`
  ```
  cmd := exec.NewCmd("", "bash", "-c", `echo 1; sleep 1; echo 1; sleep 1; echo 1; sleep 1;`)
	err := cmd.Start()
  if err != nil {
    log.Fatal(err)
  }
  for cmd.Read() {
    fmt.Println(string(cmd.GetMsg()))
  }
  err = cmd.Wait()
  if err != nil {
    log.Fatal(err)
  }
  ```
  Or setup `EventHandler`
  ```
  cmd := exec.NewCmd("", "bash", "-c", `echo 1; sleep 1; echo 1; sleep 1; echo 1; sleep 1;`)
  cmd.SetEventHandler(&exec.EventHandler{
    CmdRead: func(cmd *exec.Cmd) {
      fmt.Println(string(cmd.GetMsg()))
    },
  })
  ```
      
## Pipeline
Run set of commands. Also support `Cancel` and `STDOUT Streaming`. About `STDOUT Streaming` of pipeline, currently only support `EventHandler` way.
