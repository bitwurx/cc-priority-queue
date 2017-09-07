package task

import (
	"encoding/json"
	"time"

	"github.com/satori/go.uuid"
)

const (
	StatusPending = iota // pending task status
)

// Task is a unit of work that is queued in the priority queue
type Task struct {
	// Created is the task creation timestamp
	// Id is the unique version 1 uuid assigned for task identification
	// Key is the resource key for the task
	// Meta is user defined data that can be added to the task
	// Priority is the queue priority order
	// RunAt is a static point in time execution time
	// Status is the execution status of the task
	Created  time.Time
	Id       string
	Key      string          `json:"key"`
	Meta     json.RawMessage `json:"meta"`
	Priority float64         `json:"priority"`
	RunAt    time.Time
	Status   int
}

// NewTask returns an initialized task instance
func NewTask(data []byte) *Task {
	task := &Task{Created: time.Now(), Status: StatusPending, Id: uuid.NewV1().String()}
	json.Unmarshal(data, task)
	return task
}
