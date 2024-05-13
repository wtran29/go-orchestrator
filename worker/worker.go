package worker

import (
	"errors"
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
	Stats     *Stats                   // keep track of stats
	TaskCount int                      //keep track of number of tasks as worker
}

// RunTask handles running a task on the machine where worker is running
func (w *Worker) runTask() task.DockerResult {
	// pull task out of the queue
	t := w.Queue.Dequeue()
	if t == nil {
		log.Println("No tasks in the queue")
		return task.DockerResult{Error: nil}
	}
	// convert task from an interface to a task.Task type
	taskQueued := t.(task.Task)
	// retrieve task from worker's Db
	taskPersisted := w.Db[taskQueued.ID]
	if taskPersisted == nil {
		taskPersisted = &taskQueued
		w.Db[taskQueued.ID] = &taskQueued
	}
	var result task.DockerResult
	if task.ValidStateTransition(taskPersisted.State, taskQueued.State) {
		switch taskQueued.State {
		case task.Scheduled:
			result = w.StartTask(taskQueued)
		case task.Completed:
			result = w.StopTask(taskQueued)
		default:
			result.Error = errors.New("we should not get here")
		}
	} else {
		err := fmt.Errorf("invalid transition from %v to %v", taskPersisted.State, taskQueued.State)
		result.Error = err
		return result
	}
	return result
}

func (w *Worker) RunTasks() {
	for {
		if w.Queue.Len() != 0 {
			result := w.runTask()
			if result.Error != nil {
				log.Printf("Error running task: %v\n", result.Error)
			}
		} else {
			log.Printf("No tasks to process currently.\n")
		}
		log.Println("Sleeping for 10 seconds.")
		time.Sleep(10 * time.Second)
	}
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

func (w *Worker) GetTasks() []*task.Task {
	tasks := []*task.Task{}
	for _, t := range w.Db {
		tasks = append(tasks, t)
	}
	return tasks

}

// CollectStats periodically collect stats about the worker
func (w *Worker) CollectStats() {
	for {
		log.Println("Collecting stats")
		w.Stats = GetStats()
		w.TaskCount = w.Stats.TaskCount
		time.Sleep(15 * time.Second)
	}
}
