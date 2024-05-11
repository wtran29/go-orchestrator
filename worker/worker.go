package worker

import (
	"fmt"
	"log"
	"time"

	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
	"github.com/wtran29/go-orchestrator/task"
)

// Worker runs and keep tracks of all the tasks
type Worker struct {
	Name      string
	Queue     queue.Queue              // tasks handled in (FIFO)
	Db        map[uuid.UUID]*task.Task // to keep track of tasks
	TaskCount int                      //keep track of number of tasks as worker
}

// RunTask handles running a task on the machine where worker is running
func (w *Worker) RunTask() {
	fmt.Println("start or stop a task")
}

func (w *Worker) AddTask(t task.Task) {
	w.Queue.Enqueue(t)
}

// StartTask starts a task
func (w *Worker) StartTask(t task.Task) task.DockerResult {
	t.StartTime = time.Now().UTC()
	config := task.NewConfig(&t)
	d := task.NewDocker(config)
	result := d.Run()
	if result.Error != nil {
		log.Printf("Error running task %v: %v\n", t.ID, result.Error)
		t.State = task.Failed
		w.Db[t.ID] = &t
		return result
	}
	t.ContainerID = result.ContainerId
	t.State = task.Running
	w.Db[t.ID] = &t
	return result
}

// StopTask stops a task
func (w *Worker) StopTask(t task.Task) task.DockerResult {
	config := task.NewConfig(&t)
	d := task.NewDocker(config)

	result := d.Stop(t.ContainerID)
	if result.Error != nil {
		log.Printf("Error stopping container %v: %v\n", t.ContainerID, result.Error)
	}
	t.FinishTime = time.Now().UTC()
	t.State = task.Completed
	w.Db[t.ID] = &t
	log.Printf("Stopped and removed container %v for task %v\n", t.ContainerID, t.ID)
	return result
}

// CollectStats periodically collect stats about the worker
func (w *Worker) CollectStats() {
	fmt.Println("collect stats")
}
