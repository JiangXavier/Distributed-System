package mr

import (
	"log"
	"sync"
	"time"
)
import "net"
import "os"
import "net/rpc"
import "net/http"

type Master struct {
	// Your definitions here.
	nMap           int // number of map task
	nReduce        int // number of reduce task
	files          []string
	mapfinished    int
	reducefinished int
	maptasklog     []int // 0 not allocated;1 waiting;2 finished
	reducetasklog  []int
	mu             sync.Mutex
}

// Your code here -- RPC handlers for the worker to call.
func (m *Master) ReceiveFinishedMap(args *WorkerArgs, reply *WorkerReply) error {
	m.mu.Lock()
	m.mapfinished++
	m.maptasklog[args.MapTaskNumber] = 2
	m.mu.Unlock()
	return nil
}

func (m *Master) ReceiveFinishedReduce(args *WorkerArgs, reply *WorkerReply) error {
	m.mu.Lock()
	m.reducefinished++
	m.reducetasklog[args.ReduceTaskNumber] = 2
	m.mu.Unlock()
	return nil
}

func (m *Master) AllocateTask(args *WorkerArgs, reply *WorkerReply) error {
	m.mu.Lock()
	// if map 还没有分完
	if m.mapfinished < m.nMap {
		// allocate new map task
		// select ont to allocate
		allocate_id := -1
		for i := 0; i < m.nMap; i++ {
			if m.maptasklog[i] == 0 {
				allocate_id = i
				break
			}
		}
		if allocate_id == -1 {
			reply.TaskType = 2
			m.mu.Unlock()
		} else {
			// allocate map task
			reply.NReduce = m.nReduce
			reply.TaskType = 0
			reply.MapTaskNumber = allocate_id
			reply.Filename = m.files[allocate_id]
			m.maptasklog[allocate_id] = 1 // waiting
			m.mu.Unlock()
			// over 10s,reallocate the map task
			go func() {
				time.Sleep(time.Duration(10) * time.Second) // wait 10s
				m.mu.Lock()
				if m.maptasklog[allocate_id] == 1 {
					m.maptasklog[allocate_id] = 0
				}
				m.mu.Unlock()
			}()
		}
	} else if m.mapfinished == m.nMap && m.reducefinished < m.nReduce {
		// allocate new reduce task
		allocate_id := -1
		for i := 0; i < m.nReduce; i++ {
			if m.reducetasklog[i] == 0 {
				allocate_id = i
				break
			}
		}
		if allocate_id == -1 {
			reply.TaskType = 2
			m.mu.Unlock()
		} else {
			// allocate
			reply.NMap = m.nMap
			reply.TaskType = 1
			reply.ReduceTaskNumber = allocate_id
			m.reducetasklog[allocate_id] = 1
			m.mu.Unlock()
			go func() {
				time.Sleep(time.Duration(10) * time.Second)
				m.mu.Lock()
				if m.reducetasklog[allocate_id] == 1 {
					m.reducetasklog[allocate_id] = 0
				}
				m.mu.Unlock()
			}()
		}
	} else {
		reply.TaskType = 3
		m.mu.Unlock()
	}
	return nil
}

// an example RPC handler.
//
// the RPC argument and reply types are defined in rpc.go.
func (m *Master) Example(args *ExampleArgs, reply *ExampleReply) error {
	reply.Y = args.X + 1
	return nil
}

// start a thread that listens for RPCs from worker.go
func (m *Master) server() {
	rpc.Register(m)
	rpc.HandleHTTP()
	//l, e := net.Listen("tcp", ":1234")
	sockname := masterSock()
	os.Remove(sockname)
	l, e := net.Listen("unix", sockname)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)
}

// main/mrmaster.go calls Done() periodically to find out
// if the entire job has finished.
func (m *Master) Done() bool {
	//ret := m.mapfinished == m.reducefinished

	// Your code here.
	ret := m.reducefinished == m.nReduce
	return ret
}

// create a Master.
// main/mrmaster.go calls this function.
// nReduce is the number of reduce tasks to use.
func MakeMaster(files []string, nReduce int) *Master {
	m := Master{}

	// Your code here.
	m.files = files
	m.nMap = len(files)
	m.nReduce = nReduce
	m.maptasklog = make([]int, m.nMap)
	m.reducetasklog = make([]int, m.nReduce)
	m.server()
	return &m
}
