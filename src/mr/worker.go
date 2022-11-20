package mr

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"log"
	"net/rpc"
	"os"
	"sort"
	"time"
)

//
// Map functions return a slice of KeyValue.
//
type KeyValue struct {
	Key   string
	Value string
}

// ByKey for sorting by key.
type ByKey []KeyValue

// for sorting by key.
func (a ByKey) Len() int           { return len(a) }
func (a ByKey) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByKey) Less(i, j int) bool { return a[i].Key < a[j].Key }

//
// use ihash(key) % NReduce to choose the reduce
// task number for each KeyValue emitted by Map.
//
func ihash(key string) int {
	h := fnv.New32a()
	h.Write([]byte(key))
	return int(h.Sum32() & 0x7fffffff)
}

//
// main/mrworker.go calls this function.--
//
func Worker(mapf func(string, string) []KeyValue,
	reducef func(string, []string) string) {

	// Your worker implementation here.
	ExecuteTask(mapf, reducef)

	fmt.Println("all done")
	// uncomment to send the Example RPC to the coordinator.
	// CallExample()

}

// ExecuteTask 执行Task
func ExecuteTask(mapf func(string, string) []KeyValue,
	reducef func(string, []string) string) {
	for {
		mapReply := askForMapTask()
		if mapReply.AllMapTaskDone {
			break
		}
		doMapTask(mapReply, mapf)
	}
	for {
		reduceReply := askForReduceTask()
		if reduceReply.AllReduceTaskDone {
			return
		}
		doReduceTask(reduceReply, reducef)
	}
}

// 读取文件, 执行map task, 完成后通知coordinator
func doMapTask(reply AskMapReply, mapf func(string, string) []KeyValue) {
	if reply.MapTaskNo == -1 {
		time.Sleep(time.Second)
		return
	}
	file, err := os.Open(reply.FileName)
	if err != nil {
		log.Printf("cannot open %v", reply)
		return
	}
	content, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("cannot read %v", reply.FileName)
	}
	file.Close()
	kva := mapf(reply.FileName, string(content))
	intermediate := make(map[int][]KeyValue)
	for _, kv := range kva {
		y := ihash(kv.Key) % reply.ReduceTaskNum
		intermediate[y] = append(intermediate[y], kv)
	}
	for y, kvList := range intermediate {
		tmpFile, err := ioutil.TempFile("", "map")
		if err != nil {
			log.Fatalf("cannot create tmp file")
		}

		enc := json.NewEncoder(tmpFile)
		for _, kv := range kvList {
			err := enc.Encode(&kv)
			if err != nil {
				log.Fatalf("encode error %v", err)
			}
		}
		tmpFile.Close()
		oname := fmt.Sprintf("mr-%d-%d", reply.MapTaskNo, y)
		err = os.Rename(tmpFile.Name(), oname)
		if err != nil {
			log.Printf("map cannot rename %v", oname)
			return
		}
		log.Printf("create %v done", oname)
	}
	notifyTaskDone(0, reply.MapTaskNo)
}

// 读取中间文件, 排序, 执行reduce task, 输出为mr-out-X
func doReduceTask(reply AskReduceReply, reducef func(string, []string) string) {
	if reply.ReduceTaskNo == -1 {
		time.Sleep(time.Second)
		return
	}
	var intermediate []KeyValue
	for x := 0; x < reply.MapTaskNum; x++ {
		fileName := fmt.Sprintf("mr-%d-%d", x, reply.ReduceTaskNo)
		file, err := os.Open(fileName)
		if err != nil {
			log.Printf("cannot open %v", fileName)
			return
		}
		dec := json.NewDecoder(file)
		for {
			var kv KeyValue
			if err := dec.Decode(&kv); err != nil {
				break
			}
			intermediate = append(intermediate, kv)
		}
		file.Close()
	}
	sort.Sort(ByKey(intermediate))
	oname := fmt.Sprintf("mr-out-%d", reply.ReduceTaskNo)
	tmpFile, err := ioutil.TempFile("", "reduce")
	if err != nil {

	}
	i := 0
	for i < len(intermediate) {
		j := i + 1
		for j < len(intermediate) && intermediate[j].Key == intermediate[i].Key {
			j++
		}
		values := []string{}
		// 把相同key的所有value放到一起
		for k := i; k < j; k++ {
			values = append(values, intermediate[k].Value)
		}
		output := reducef(intermediate[i].Key, values)

		// this is the correct format for each line of Reduce output.
		fmt.Fprintf(tmpFile, "%v %v\n", intermediate[i].Key, output)

		i = j
	}
	tmpFile.Close()
	err = os.Rename(tmpFile.Name(), oname)
	if err != nil {
		log.Printf("reduce cannot rename %v", oname)
		return
	}
	notifyTaskDone(1, reply.ReduceTaskNo)
}

//
// example function to show how to make an RPC call to the coordinator.
//
// the RPC argument and reply types are defined in rpc.go.
//
func CallExample() {

	// declare an argument structure.
	args := ExampleArgs{}

	// fill in the argument(s).
	args.X = 99

	// declare a reply structure.
	reply := ExampleReply{}

	// send the RPC request, wait for the reply.
	// the "Coordinator.Example" tells the
	// receiving server that we'd like to call
	// the Example() method of struct Coordinator.
	ok := call("Coordinator.Example", &args, &reply)
	if ok {
		// reply.Y should be 100.
		fmt.Printf("reply.Y %v\n", reply.Y)
	} else {
		fmt.Printf("call failed!\n")
	}
}

func askForMapTask() AskMapReply {
	args := AskMapReq{}
	reply := AskMapReply{}
	ok := call("Coordinator.MapTask", args, &reply)
	if ok {
		fmt.Printf("askForMapTask %v\n", reply)
	} else {
		fmt.Printf("call failed!\n")
	}
	return reply
}

func askForReduceTask() AskReduceReply {
	args := AskReduceReq{}
	reply := AskReduceReply{}
	ok := call("Coordinator.ReduceTask", args, &reply)
	if ok {
		fmt.Printf("askForReduceTask %v\n", reply)
	} else {
		fmt.Printf("call failed!\n")
	}
	return reply
}

func notifyTaskDone(taskType, taskNo int) {
	args := NotifyTaskDoneReq{
		TaskType: taskType,
		TaskNo:   taskNo,
	}
	reply := NotifyTaskDoneReply{}
	ok := call("Coordinator.NotifyTaskDone", args, &reply)
	if ok {
		fmt.Printf("notifyTaskDone %v\n", args)
	} else {
		fmt.Printf("call failed!\n")
	}
}

//
// send an RPC request to the coordinator, wait for the response.
// usually returns true.
// returns false if something goes wrong.
//
func call(rpcname string, args interface{}, reply interface{}) bool {
	// c, err := rpc.DialHTTP("tcp", "127.0.0.1"+":1234")
	sockname := coordinatorSock()
	c, err := rpc.DialHTTP("unix", sockname)
	if err != nil {
		log.Fatal("dialing:", err)
	}
	defer c.Close()

	err = c.Call(rpcname, args, reply)
	if err == nil {
		return true
	}

	fmt.Println(err)
	return false
}
