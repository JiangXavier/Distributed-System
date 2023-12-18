package mr

//
// RPC definitions.
//
// remember to capitalize all names.
//

import (
	"os"
)
import "strconv"

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
type WorkerArgs struct {
	MapTaskNumber    int // id of finished maptask
	ReduceTaskNumber int // id of finished reducetask
}
type WorkerReply struct {
	TaskType         int    // 0 : map;1 : reduce; 2:waiting ; 3:finished
	NMap             int    // map task number
	NReduce          int    // reduce task number
	MapTaskNumber    int    // id of maptask only
	ReduceTaskNumber int    // id of reducetask
	Filename         string //map only
}

// Cook up a unique-ish UNIX-domain socket name
// in /var/tmp, for the master.
// Can't use the current directory since
// Athena AFS doesn't support UNIX-domain sockets.
func masterSock() string {
	s := "/var/tmp/824-mr-"
	s += strconv.Itoa(os.Getuid())
	return s
}
