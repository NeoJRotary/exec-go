package exec

// import "sync"

// // PipelineGroup ...
// type PipelineGroup struct {
// 	ex   *Exec
// 	Stop bool
// 	List []*Pipeline
// }

// var mutex = &sync.Mutex{}

// // NewPipelineGroup ...
// func (ex *Exec) NewPipelineGroup() *PipelineGroup {
// 	group := PipelineGroup{
// 		ex:   ex,
// 		Stop: false,
// 		List: []*Pipeline{},
// 		// cmds:     []*Cmd{},
// 	}
// 	return &group
// }

// // NewPipeline ...
// func (group *PipelineGroup) NewPipeline() *Pipeline {
// 	pipe := group.ex.NewPipeline()
// 	pipe.group = group
// 	mutex.Lock()
// 	group.List = append(group.List, pipe)
// 	mutex.Unlock()
// 	return pipe
// }

// // Cancel ...
// func (group *PipelineGroup) Cancel() {
// 	for _, p := range group.List {
// 		p.Cancel()
// 	}
// 	group.Stop = true
// }

// func (group *PipelineGroup) childStop() {
// 	for _, p := range group.List {
// 		if !p.Stop {
// 			return
// 		}
// 	}
// 	group.Stop = true
// }
