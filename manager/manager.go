package manager

import (
	"fmt"

	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
	"github.com/wtran29/go-orchestrator/task"
)

// Manager will keep track of the workers in the cluster
type Manager struct {
	Pending       queue.Queue // which tasks will be placed upon first being submitted
	TaskDb        map[string][]task.Task
	EventDb       map[string][]task.TaskEvent
	Workers       []string // keep track of the workers
	WorkerTaskMap map[string][]uuid.UUID
	TaskWorkerMap map[uuid.UUID]string
}

// SelectWorker is responsible for look at requirements specified in a Task
// and evaluating the resources available in pool of workers to pick the right
// worker for the task
func (m *Manager) SelectWorker() {
	fmt.Println("select worker")
}

// UpdateTasks will track tasks, their states, and machine on which they run
func (m *Manager) UpdateTasks() {
	fmt.Println("update tasks")
}

// SendWork sends tasks to workers
func (m *Manager) SendWork() {
	fmt.Println("send work to workers")
}
