package mr

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"sync"
	"time"
)

const (
	// DefaultTimeout 10s内要得知woker是否完工, 超时则把相同的task发给另一个worker
	DefaultTimeout = 10 * time.Second
)

type Coordinator struct {
	// Your definitions here.
	mutex  sync.Mutex
	mtasks []MapTask
	rtasks []ReduceTask
	done   chan bool
}

// MapTask startTime为nil则未调度, startTime不为nil,endTime为nil为执行中, endTime为nil则执行完成
type MapTask struct {
	startTime *time.Time
	endTime   *time.Time
	fileName  string
}

type ReduceTask struct {
	startTime *time.Time
	endTime   *time.Time
}

func (c *Coordinator) init(files []string, nReduce int) {
	c.mtasks = make([]MapTask, len(files))
	for i, file := range files {
		c.mtasks[i].fileName = file
	}
	c.rtasks = make([]ReduceTask, nReduce)
	c.done = make(chan bool, 1)
}

// 每过1s看下哪些task超时了, 设置成未执行
func (c *Coordinator) polling() {
	t := time.NewTicker(time.Second)
	defer t.Stop()
	for {
		select {
		case <-t.C:
		case <-c.done:
			return
		default:
			c.mutex.Lock()
			scanMapTasks(c.mtasks)
			scanReduceTasks(c.rtasks)
			c.mutex.Unlock()
		}
	}
}

func scanMapTasks(mtasks []MapTask) {
	for i := range mtasks {
		if mtasks[i].startTime == nil || mtasks[i].endTime != nil {
			continue
		}
		if time.Until(*mtasks[i].startTime) > DefaultTimeout {
			mtasks[i].startTime = nil
		}
	}
}

func scanReduceTasks(rtasks []ReduceTask) {
	for i := range rtasks {
		if rtasks[i].startTime == nil || rtasks[i].endTime != nil {
			continue
		}
		if time.Until(*rtasks[i].startTime) > DefaultTimeout {
			rtasks[i].startTime = nil
		}
	}
}

// Your code here -- RPC handlers for the worker to call.

// MapTask 回复worker要执行哪个map task
func (c *Coordinator) MapTask(args AskMapReq, reply *AskMapReply) error {
	allDone := true
	c.mutex.Lock()
	defer c.mutex.Unlock()
	for i := range c.mtasks {
		if c.mtasks[i].startTime == nil {
			curTime := time.Now()
			c.mtasks[i].startTime = &curTime
			reply.MapTaskNo = i
			reply.ReduceTaskNum = len(c.rtasks)
			reply.FileName = c.mtasks[i].fileName
			return nil
		}
		if allDone && c.mtasks[i].endTime == nil {
			allDone = false
		}
	}
	reply.MapTaskNo = -1
	reply.AllMapTaskDone = allDone
	return nil
}

// ReduceTask 回复worker要执行哪个reduce task
func (c *Coordinator) ReduceTask(args AskReduceReq, reply *AskReduceReply) error {
	allDone := true
	c.mutex.Lock()
	defer c.mutex.Unlock()
	for i := range c.rtasks {
		if c.rtasks[i].startTime == nil {
			curTime := time.Now()
			c.rtasks[i].startTime = &curTime
			reply.ReduceTaskNo = i
			reply.MapTaskNum = len(c.mtasks)
			return nil
		}
		if allDone && c.rtasks[i].endTime == nil {
			allDone = false
		}
	}
	reply.ReduceTaskNo = -1 // 无可以提供的task(比如都在运行)
	reply.AllReduceTaskDone = allDone
	return nil
}

// NotifyTaskDone worker告知哪个task完成了
func (c *Coordinator) NotifyTaskDone(args NotifyTaskDoneReq, reply *NotifyTaskDoneReply) error {
	curTime := time.Now()
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if args.TaskType == 0 && c.mtasks[args.TaskNo].endTime == nil {
		c.mtasks[args.TaskNo].endTime = &curTime
	} else if c.rtasks[args.TaskNo].endTime == nil {
		c.rtasks[args.TaskNo].endTime = &curTime
	}
	return nil
}

//
// an example RPC handler.
//
// the RPC argument and reply types are defined in rpc.go.
//
func (c *Coordinator) Example(args *ExampleArgs, reply *ExampleReply) error {
	reply.Y = args.X + 1
	return nil
}

//
// start a thread that listens for RPCs from worker.go
//
func (c *Coordinator) server() {
	rpc.Register(c)
	rpc.HandleHTTP()
	//l, e := net.Listen("tcp", ":1234")
	sockname := coordinatorSock()
	os.Remove(sockname)
	l, e := net.Listen("unix", sockname)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)
}

//
// main/mrcoordinator.go calls Done() periodically to find out
// if the entire job has finished.
//
func (c *Coordinator) Done() bool {
	// Your code here.
	c.mutex.Lock()
	defer c.mutex.Unlock()
	for i := range c.rtasks {
		if c.rtasks[i].endTime == nil {
			return false
		}
	}
	c.done <- true
	return true
}

//
// create a Coordinator.
// main/mrcoordinator.go calls this function.
// nReduce is the number of reduce tasks to use.
//
func MakeCoordinator(files []string, nReduce int) *Coordinator {
	c := Coordinator{}

	// Your code here.
	c.init(files, nReduce)
	go c.polling()
	c.server()
	return &c
}
