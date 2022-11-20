package mr

//
// RPC definitions.
//
// remember to capitalize all names.
//

import (
	"os"
	"strconv"
)

//
// example to show how to declare the arguments
// and reply for an RPC.
//

type ExampleArgs struct {
	X int
}

type ExampleReply struct {
	Y int
}

// Add your RPC definitions here.
type AskMapReq struct {
}

type AskMapReply struct {
	MapTaskNo      int    // 0 ~ m-1 个map Task的哪一个
	ReduceTaskNum  int    // reduce阶段有几个task
	FileName       string // 哪个input File
	AllMapTaskDone bool   // 如果所有map task都已完成, 才可以做Reduce task
}

type AskReduceReq struct {
}

type AskReduceReply struct {
	MapTaskNum        int // map阶段有几个task, 方便遍历所需的intermediate files
	ReduceTaskNo      int // 0 ~ r-1 个reduce Task的哪一个
	AllReduceTaskDone bool
}

// NotifyTaskDoneReq 告知reduce task已完成
type NotifyTaskDoneReq struct {
	TaskType int // 0: map task, 1: reduce task
	TaskNo   int // 哪个task完成了
}

type NotifyTaskDoneReply struct {
}

// Cook up a unique-ish UNIX-domain socket name
// in /var/tmp, for the coordinator.
// Can't use the current directory since
// Athena AFS doesn't support UNIX-domain sockets.
func coordinatorSock() string {
	s := "/var/tmp/824-mr-"
	s += strconv.Itoa(os.Getuid())
	return s
}
