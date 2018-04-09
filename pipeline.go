package exec

// import (
// 	"fmt"
// 	"log"
// )

// // Pipeline commands pipeline
// type Pipeline struct {
// 	Started  bool
// 	Done     bool
// 	Failed   bool
// 	Canceled bool
// 	Pausing  bool
// 	logger   *log.Logger
// 	stopChan chan bool
// 	group    *PipelineGroup
// 	underDir string
// 	pipeCmds []*pipeCmd
// 	index    int
// }

// // HandleFunc handler between two commands in pipeline.
// type HandleFunc func(string, error) bool

// type pipeCmd struct {
// 	cmd     *Cmd
// 	handler HandleFunc
// }

// // Pipeline HandleFunc to stop it at error
// func pipeStopAtErr(s string, e error) bool {
// 	return e == nil
// }

// // Pipeline HandleFunc to skip error
// func pipeSkipErr(string, error) bool {
// 	return true
// }

// // NewPipeline new pipeline.
// func (ex *Exec) NewPipeline(withLogger bool) *Pipeline {
// 	pipe := Pipeline{
// 		stopChan: make(chan bool, 1),
// 		underDir: "",
// 		index:    0,
// 	}
// 	if withLogger {
// 		pipe.logger = ex.logger
// 	}
// 	return &pipe
// }

// // Under under which dir. Return pipeline pointer.
// func (pipe *Pipeline) Under(dir string) *Pipeline {
// 	if pipe.Started {
// 		return pipe
// 	}
// 	pipe.underDir = dir
// 	return pipe
// }

// // Do add cmd to pipeline. It use last dir which you set by Under(). By default dir is current path.
// func (pipe *Pipeline) Do(name string, arg ...string) *Pipeline {
// 	if pipe.Started {
// 		return pipe
// 	}
// 	cmd := &pipeCmd{
// 		cmd:     NewCmd(pipe.underDir, name, arg...),
// 		handler: pipeStopAtErr,
// 	}
// 	pipe.pipeCmds = append(pipe.pipeCmds, cmd)
// 	return pipe
// }

// func (pipe *Pipeline) lastPipeCmd() *pipeCmd {
// 	return pipe.pipeCmds[len(pipe.pipeCmds)-1]
// }

// // WithEnv under which dir. Return pipeline pointer.
// func (pipe *Pipeline) WithEnv(name, value string) *Pipeline {
// 	if pipe.Started {
// 		return pipe
// 	}
// 	pipe.lastPipeCmd().cmd.AddEnv(name, value)
// 	return pipe
// }

// // Handler add handler after cmd for handling last output and error.
// func (pipe *Pipeline) Handler(f HandleFunc) *Pipeline {
// 	if pipe.Started {
// 		return pipe
// 	}
// 	pipe.lastPipeCmd().handler = f
// 	return pipe
// }

// // SkipErr add handler after cmd to skip err.
// func (pipe *Pipeline) SkipErr(f HandleFunc) *Pipeline {
// 	if pipe.Started {
// 		return pipe
// 	}
// 	pipe.lastPipeCmd().handler = pipeSkipErr
// 	return pipe
// }

// // Start start pipeline
// func (pipe *Pipeline) Start() {
// 	pipe.Started = true
// 	go pipe.run()
// }

// func (pipe *Pipeline) run() {
// 	for _, c := range pipe.pipeCmds {

// 	}
// }

// func (pipe *Pipeline) stop() {
// 	pipe.Canceled = true
// 	if pipe.group != nil {
// 		pipe.group.childStop()
// 	}
// }

// // Cancel ...
// func (pipe *Pipeline) Cancel() {
// 	if pipe.Canceled {
// 		return
// 	}
// 	pipe.stopChan <- true
// 	pipe.stop()
// }

// // Wait ...
// func (pipe *Pipeline) Wait() (result []Result, cancel bool, err error) {
// 	if pipe.Canceled {
// 		return nil, true, nil
// 	}
// 	for _, c := range pipe.cmds {
// 		if pipe.ex.doStdLog {
// 			fmt.Println(c.str)
// 		}
// 		c.Start()
// 		go c.Wait()
// 		select {
// 		case <-pipe.stopChan:
// 			c.Cancel()
// 			return nil, true, nil
// 		case err = <-c.done:
// 			result = append(result, Result{
// 				cmd:    c.str,
// 				output: c.stdout.String(),
// 			})
// 			if pipe.ex.doStdLog {
// 				if err != nil {
// 					pipe.ex.stdLogger.Println(c.stdout.String(), err.Error())
// 				} else {
// 					pipe.ex.stdLogger.Println(c.stdout.String())
// 				}
// 			}
// 			if err != nil {
// 				return result, false, err
// 			}
// 		}
// 	}
// 	pipe.stop()
// 	return result, false, nil
// }
