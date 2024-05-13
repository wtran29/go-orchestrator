package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/docker/docker/client"
	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/wtran29/go-orchestrator/manager"
	"github.com/wtran29/go-orchestrator/task"
	"github.com/wtran29/go-orchestrator/worker"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic("Error loading .env file")
	}
	host := os.Getenv("ARCHON_HOST")
	port, _ := strconv.Atoi(os.Getenv("ARCHON_PORT"))

	db := make(map[uuid.UUID]*task.Task)

	fmt.Println("Starting Archon worker")

	w := worker.Worker{
		Queue: *queue.New(),
		Db:    db,
	}
	api := worker.Api{Address: host, Port: port, Worker: &w}

	go runTasks(&w)
	go w.CollectStats()
	go api.Start()

	workers := []string{fmt.Sprintf("%s:%d", host, port)}
	m := manager.New(workers)

	for i := 0; i < 3; i++ {
		t := task.Task {
			ID: uuid.New(),
			Name: fmt.Sprintf("test-container-%d", i),
			State: task.Scheduled,
			Image: "strm/helloworld-http",
		}
		te := task.TaskEvent{
			ID: uuid.New(),
			State: task.Running,
			Task: t,
		}
		m.AddTask(te)
		m.SendWork()
	}
	go func () {
		for {
			fmt.Printf("[Manager] Updating tasks from %d workers\n", len(workers))
			m.UpdateTasks()
			time.Sleep(15 * time.Second)
		}
	}()

	for {
		for _, t := range m.TaskDb {
			fmt.Printf("[Manager] Tasks: id %s, state: %d\n", t.ID, t.State)
			time.Sleep(15 * time.Second)
		}
	}
}

func createContainer() (*task.Docker, *task.DockerResult) {
	c := task.Config{
		Name:  "test-container-1",
		Image: "postgres:13",
		Env: []string{
			"POSTGRES_USER=archon",
			"POSTGRES_PASSWORD=secret",
		},
	}
	dc, _ := client.NewClientWithOpts(client.FromEnv)
	d := task.Docker{
		Client: dc,
		Config: c,
	}
	result := d.Run()
	if result.Error != nil {
		fmt.Printf("%v\n", result.Error)
		return nil, nil
	}
	fmt.Printf("Container %s is running with config %v\n", result.ContainerId, c)
	return &d, &result
}

func stopContainer(d *task.Docker, id string) *task.DockerResult {
	result := d.Stop(id)
	if result.Error != nil {
		fmt.Printf("%v\n", result.Error)
	}
	fmt.Printf("Container %s has been stopped and removed\n", id)
	return &result
}

func runTasks(w *worker.Worker) {
	for {
		if w.Queue.Len() != 0 {
			result := w.RunTask()
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
