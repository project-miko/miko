package taskpool

import (
	"sync"

	"github.com/project-miko/miko/tools/log"
)

const (
	taskWorkerNum = 10
)

var (
	// task worker channel
	workerChs [taskWorkerNum]chan *Task
	// load balancing: average strategy
	workBalanceCounter = 0
	workBalanceLocker  sync.Mutex
)

type Task struct {
	CbFunc func(msg interface{})
	Msg    interface{}
}

// task load balancing atomic counter
func getBalanceNum() int {
	workBalanceLocker.Lock()
	defer workBalanceLocker.Unlock()
	workBalanceCounter++
	if workBalanceCounter >= taskWorkerNum {
		workBalanceCounter = 0
	}
	return workBalanceCounter
}

func AddPushTask(t *Task) {
	wid := getBalanceNum()
	workerChs[wid] <- t
}

// initialize task worker
func InitTaskListeners() {

	for i := 0; i < taskWorkerNum; i++ {

		go func(index int) {
			defer log.PrintPanicStackError()
			workerChs[index] = make(chan *Task, 10000)

			for {
				select {
				case task := <-workerChs[index]:
					l := len(workerChs[index])
					if l > 5000 {
						log.Warning("", "task pool channel(%d) overstock (greater then 5k)", l)
					}
					// discard if greater than 8000
					if l > 8000 {
						log.Error("", "task pool channel(%d) overstock (greater then 8k)", l)
						continue
					}

					task.CbFunc(task.Msg)
				}
			}
		}(i)
	}
}
