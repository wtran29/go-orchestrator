package task

import (
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/google/uuid"
)

type State int

const (
	Pending State = iota
	Scheduled
	Running
	Completed
	Failed
)

// Task represents a task that auser wants to run on cluster
type Task struct {
	ID            uuid.UUID
	Name          string
	State         State
	Image         string // what docker image task should use
	Memory        int // amount of memory needed
	Disk          int // amount of disk space needed
	ExposedPorts  nat.PortSet
	PortBindings  map[string]string
	RestartPolicy string
	StartTime     time.Time
	FinishTime    time.Time
}

// TaskEvent represents an even that moves a Task from
// one state to another
type TaskEvent struct {
	ID uuid.UUID
	State State
	Timestamp time.Time
	Task Task
}