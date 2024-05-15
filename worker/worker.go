package worker

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/golang-collections/collections/queue"
	"github.com/wtran29/go-orchestrator/stats"
	"github.com/wtran29/go-orchestrator/store"
	"github.com/wtran29/go-orchestrator/task"
)

// Worker runs and keep tracks of all the tasks
type Worker struct {
	Name      string
	Queue     queue.Queue  // tasks handled in (FIFO)
	Db        store.Store  // to keep track of tasks
	Stats     *stats.Stats // keep track of stats
	TaskCount int          //keep track of number of tasks as worker
}

func New(name string, taskDBtype string) *Worker {
	w := Worker{
		Name:  name,
		Queue: *queue.New(),
	}
	var s store.Store
	var err error
	switch taskDBtype {
	case "memory":
		s = store.NewInMemoryTaskStore()
	case "persistent":
		filename := fmt.Sprintf("%s_task.db", name)
		s, err = store.NewTaskStore(filename, 0600, "tasks")
	}
	if err != nil {
		log.Printf("unable to create new task store: %v", err)
	}
	w.Db = s
	return &w

}

// RunTask handles running a task on the machine where worker is running
func (w *Worker) runTask() task.DockerResult {
	// pull task out of the queue
	t := w.Queue.Dequeue()
	if t == nil {
		log.Println("[worker] No tasks in the queue")
		return task.DockerResult{Error: nil}
	}
	// convert task from an interface to a task.Task type
	taskQueued := t.(task.Task)
	fmt.Printf("[worker] Found task in queue: %v\n", taskQueued)

	err := w.Db.Put(taskQueued.ID.String(), &taskQueued)
	if err != nil {
		msg := fmt.Errorf("error stroing task %s: %v", taskQueued.ID.String(), err)
		log.Println(msg)
		return task.DockerResult{Error: msg}
	}
	result, err := w.Db.Get(taskQueued.ID.String())
	if err != nil {
		msg := fmt.Errorf("error getting task %s from database: %v", taskQueued.ID.String(), err)
		log.Println(msg)
		return task.DockerResult{Error: msg}
	}
	// retrieve task from worker's Db
	taskPersisted := *result.(*task.Task)
	if taskPersisted.State == task.Completed {
		return w.StopTask(taskPersisted)
	}
	var dockerResult task.DockerResult
	if task.ValidStateTransition(taskPersisted.State, taskQueued.State) {
		switch taskQueued.State {
		case task.Scheduled:
			if taskQueued.ContainerID != "" {
				dockerResult = w.StopTask(taskQueued)
				if dockerResult.Error != nil {
					log.Printf("%v\n", dockerResult.Error)
				}
			}
			dockerResult = w.StartTask(taskQueued)
		default:
			fmt.Printf("This is a mistake. taskPersisted: %v, taskQueued: %v\n", taskPersisted, taskQueued)
			dockerResult.Error = errors.New("we should not get here")
		}
	} else {
		err := fmt.Errorf("invalid transition from %v to %v", taskPersisted.State, taskQueued.State)
		dockerResult.Error = err
		return dockerResult
	}
	return dockerResult
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
	config := task.NewConfig(&t)
	d := task.NewDocker(config)
	result := d.Run()
	if result.Error != nil {
		log.Printf("Error running task %v: %v\n", t.ID, result.Error)
		t.State = task.Failed
		w.Db.Put(t.ID.String(), &t)
		return result
	}
	t.ContainerID = result.ContainerId
	t.State = task.Running
	w.Db.Put(t.ID.String(), &t)
	return result
}

// StopTask stops a task
func (w *Worker) StopTask(t task.Task) task.DockerResult {
	config := task.NewConfig(&t)
	d := task.NewDocker(config)

	stopResult := d.Stop(t.ContainerID)
	if stopResult.Error != nil {
		log.Printf("Error stopping container: %v\n", stopResult.Error)
	}
	removeResult := d.Remove(t.ContainerID)
	if removeResult.Error != nil {
		log.Printf("Error removing container: %v\n", removeResult.Error)
	}
	t.FinishTime = time.Now().UTC()
	t.State = task.Completed
	w.Db.Put(t.ID.String(), &t)
	log.Printf("Stopped and removed container %v for task %v\n", t.ContainerID, t.ID)
	return removeResult
}

func (w *Worker) GetTasks() []*task.Task {
	taskList, err := w.Db.List()
	if err != nil {
		log.Printf("error getting list of tasks: %v", err)
		return nil
	}

	return taskList.([]*task.Task)

}

// CollectStats periodically collect stats about the worker
func (w *Worker) CollectStats() {
	for {
		log.Println("Collecting stats")
		w.Stats = stats.GetStats()
		w.TaskCount = w.Stats.TaskCount
		time.Sleep(15 * time.Second)
	}
}

// InspectTask method calls new Inspect method on the Docker struct
func (w *Worker) InspectTask(t task.Task) task.DockerInspectResponse {
	config := task.NewConfig(&t)
	d := task.NewDocker(config)
	return d.Inspect(t.ContainerID)
}

// UpdateTasks serves as a wrapper to updateTasks method
func (w *Worker) UpdateTasks() {
	for {
		log.Println("Checking status of tasks")
		w.updateTasks()
		log.Println("Task updates completed")
		log.Println("Sleeping for 15 seconds")
		time.Sleep(15 * time.Second)
	}
}

func (w *Worker) updateTasks() {
	tasks, err := w.Db.List()
	if err != nil {
		log.Printf("error getting list of tasks: %v\n", err)
		return
	}
	for _, t := range tasks.([]*task.Task) {
		if t.State == task.Running {
			// call InspectTask method to get task state from docker daemon
			resp := w.InspectTask(*t)
			if resp.Error != nil {
				fmt.Printf("ERROR: %v\n", resp.Error)
			}
			// verify task is in running state
			if resp.Container == nil {
				log.Printf("No container for running task %s\n", t.ID)
				t.State = task.Failed
				w.Db.Put(t.ID.String(), t)
			}
			// if not in running state, or not running at all, set task state to failed
			if resp.Container.State.Status == "exited" {
				log.Printf("Container for task %s in non-running state %s\n", t.ID, resp.Container.State.Status)
				t.State = task.Failed
				w.Db.Put(t.ID.String(), t)
			}
			// task is running, update exposed ports
			t.HostPorts = resp.Container.NetworkSettings.NetworkSettingsBase.Ports
			w.Db.Put(t.ID.String(), t)
		}
	}
}
