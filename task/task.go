package task

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/go-connections/nat"
	"github.com/google/uuid"
)

// Task represents a task that a user wants to run on cluster
type Task struct {
	ID            uuid.UUID
	ContainerID   string
	Name          string
	State         State
	Image         string  // what docker image task should use
	Cpu           float64 // amount of cpu usage
	Memory        int64   // amount of memory needed
	Disk          int64   // amount of disk space needed
	ExposedPorts  nat.PortSet
	PortBindings  map[string]string
	RestartPolicy string // ["", "always", "unless-stopped", "on-failure"]
	StartTime     time.Time
	FinishTime    time.Time
}

// TaskEvent represents an even that moves a Task from
// one state to another
type TaskEvent struct {
	ID        uuid.UUID
	State     State
	Timestamp time.Time
	Task      Task
}

// Config struct to hold Docker container config
type Config struct {
	Name         string // name of the container
	AttachStdin  bool
	AttachStdout bool
	Attachstderr bool
	ExposedPorts nat.PortSet // list of ports exposed
	Cmd          []string
	Image        string // Image used to run the container
	// Memory and Disk serves two purposes:
	// scheduler use them to find node in cluster
	Memory        int64    // Memory in MiB
	Disk          int64    // Disk in GiB
	Env           []string // allows user to specify env variables passed into the container
	RestartPolicy string   // tells Docker daemon what to do in event container dies
}

func NewConfig(t *Task) *Config {
	return &Config{
		Name:          t.Name,
		ExposedPorts:  t.ExposedPorts,
		Image:         t.Image,
		Memory:        t.Memory,
		Disk:          t.Disk,
		RestartPolicy: t.RestartPolicy,
	}
}

// Docker represents the Docker container
type Docker struct {
	Client *client.Client // Docker client object
	Config Config         //holds the task configuration
}

func NewDocker(c *Config) *Docker {
	dc, _ := client.NewClientWithOpts(client.FromEnv)
	return &Docker{Client: dc, Config: *c}
}

// DockerResult represents the Docker results
type DockerResult struct {
	Error       error
	Action      string
	ContainerId string
	Result      string
}

// Run pulls the container's image
func (d *Docker) Run() DockerResult {
	ctx := context.Background()
	reader, err := d.Client.ImagePull(ctx, d.Config.Image, types.ImagePullOptions{})
	if err != nil {
		log.Printf("Error pulling image %s: %v\n", d.Config.Image, err)
		return DockerResult{Error: err}
	}
	io.Copy(os.Stdout, reader)

	rp := container.RestartPolicy{
		Name: d.Config.RestartPolicy,
	}

	r := container.Resources{
		Memory: d.Config.Memory,
	}
	cc := container.Config{
		Image:        d.Config.Image,
		Tty:          false,
		Env:          d.Config.Env,
		ExposedPorts: d.Config.ExposedPorts,
	}
	hc := container.HostConfig{
		RestartPolicy:   rp,
		Resources:       r,
		PublishAllPorts: true,
	}
	resp, err := d.Client.ContainerCreate(ctx, &cc, &hc, nil, nil, d.Config.Name)
	if err != nil {
		log.Printf("Error creating container using image %s: %v\n", d.Config.Image, err)
		return DockerResult{Error: err}
	}
	err = d.Client.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
	if err != nil {
		log.Printf("Error starting container %s: %v\n", resp.ID, err)
		return DockerResult{Error: err}
	}

	out, err := d.Client.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
	})
	if err != nil {
		log.Printf("Error getting logs for container %s: %v\n", resp.ID, err)
	}

	stdcopy.StdCopy(os.Stdout, os.Stderr, out)

	return DockerResult{ContainerId: resp.ID, Action: "start", Result: "success"}

}

func (d *Docker) Stop(id string) DockerResult {
	log.Printf("Attempting to stop container %v", id)
	ctx := context.Background()
	err := d.Client.ContainerStop(ctx, id, container.StopOptions{})
	if err != nil {
		fmt.Printf("Error stopping container %s: %v\n", id, err)
		return DockerResult{Error: err}
	}
	err = d.Client.ContainerRemove(ctx, id, types.ContainerRemoveOptions{
		RemoveVolumes: true,
		RemoveLinks:   false,
		Force:         false,
	})
	if err != nil {
		log.Printf("Error removing container %s: %v\n", id, err)
		return DockerResult{}
	}
	return DockerResult{Action: "stop", Result: "success", Error: nil}
}
