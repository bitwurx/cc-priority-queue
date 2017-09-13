package tasks

import (
	"encoding/json"
	"time"

	arango "github.com/arangodb/go-driver"
	"github.com/satori/go.uuid"

	"concord-pq/database"
)

const (
	StatusPending = iota // pending task status.
	StatusQueued         // queued task status.
)

const (
	CollectionTasks     = "tasks"      // the name of the tasks database collection
	CollectionTaskStats = "task_stats" // the name of the task stats database collection
)

// TaskStat stores a runtime for a task.
type TaskStat struct {
	// Created is the task stat creation timestamp.
	// Key is the task key.
	// Runtime is the task run time in seconds.
	Created time.Time `json:"created"`
	Key     string    `json:"key"`
	Runtime float64   `json:"runtime"`
}

// NewTaskStat returns an initialized task instance.
func NewTaskStat(key string, runtime float64) *TaskStat {
	return &TaskStat{time.Now(), key, runtime}
}

// Save creates a new document for the task stat in the database.
func (taskStat *TaskStat) Save(m database.Model) (arango.DocumentMeta, error) {
	return m.Save(taskStat)
}

// Task is a unit of work that is queued in the priority queue.
type Task struct {
	// Created is the task creation timestamp.
	// Id is the unique version 1 uuid assigned for task identification.
	// Key is the resource key for the task.
	// Meta is user defined data that can be added to the task.
	// Priority is the queue priority order.
	// RunAt is a static point in time execution time.
	// Status is the execution status of the task.
	Created  time.Time       `json:"created"`
	Id       string          `json:"_key"`
	Key      string          `json:"key"`
	Meta     json.RawMessage `json:"meta;omitempty"`
	Priority float64         `json:"priority"`
	RunAt    time.Time       `json:"runAt;omitempty"`
	Status   int             `json:"status"`
}

// NewTask returns an initialized task instance.
func NewTask(data []byte) *Task {
	task := &Task{Created: time.Now(), Status: StatusPending, Id: uuid.NewV1().String()}
	json.Unmarshal(data, task)
	return task
}

// ChangeStatus changes the status of the task and saves the task
func (task *Task) ChangeStatus(model database.Model, status int) error {
	task.Status = status
	_, err := model.Save(task)
	return err
}

// Save writes the task to the database
func (task *Task) Save(model database.Model) (arango.DocumentMeta, error) {
	return model.Save(task)
}
