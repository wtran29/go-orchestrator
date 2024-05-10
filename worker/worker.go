package worker

import (
	"fmt"

	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
	"github.com/wtran29/go-orchestrator/task"
)

// Worker runs and keep tracks of all the tasks
type Worker struct {
	Name      string
	Queue     queue.Queue             // tasks handled in (FIFO)
	Db        map[uuid.UUID]task.Task // to keep track of tasks
	TaskCount int                     //keep track of number of tasks as worker
}

// RunTask handles running a task on the machine where worker is running
func (w *Worker) RunTask() {
	fmt.Println("start or stop a task")
}

// StartTask starts a task
func (w *Worker) StartTask() {
	fmt.Println("start a task")
}

// StopTask stops a task
func (w *Worker) StopTask() {
	fmt.Println("stop a task")
}

// CollectStats periodically collect stats about the worker
func (w *Worker) CollectStats() {
	fmt.Println("collect stats")
}
