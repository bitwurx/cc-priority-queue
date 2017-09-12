package task

import (
	"encoding/json"
	"time"

	arango "github.com/arangodb/go-driver"
	"github.com/satori/go.uuid"

	. "concord-pq/database"
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

func (taskStat TaskStat) CreateCollection(db Database) {
	db.CreateCollection(nil, CollectionTaskStats, nil)
}

// Save creates a new document for the task stat in the database.
func (taskStat *TaskStat) Save(taskStats Collection) {
	taskStats.CreateDocument(nil, taskStat)
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

// ChangeStatus changes the status of the task.
//
// The task is saved to the database after status change.
func (task *Task) ChangeStatus(tasks Collection, status int) error {
	task.Status = status
	return task.Save(tasks)
}

func (task Task) CreateCollection(db Database) {
	db.CreateCollection(nil, CollectionTasks, nil)
}

// Save creates a new document for the task or updates the existing task status
// if a matching task id exists.
func (task *Task) Save(tasks Collection) error {
	_, err := tasks.CreateDocument(nil, task)
	if err != nil && arango.IsConflict(err) {
		patch := map[string]interface{}{"Status": task.Status}
		_, err = tasks.UpdateDocument(nil, task.Id, patch)
	}
	return nil
}

// Initialize database models
func init() {
	AddModel(Task{})
	AddModel(TaskStat{})
}
