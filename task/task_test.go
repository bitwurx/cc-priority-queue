package task_test

import (
	"testing"
	"time"

	"github.com/satori/go.uuid"

	"concord-pq/task"
)

func TestNewTask(t *testing.T) {
	data := `{"meta": {"id": 123}, "Priority": 22.5, "key": "tb1"}`
	task := task.NewTask([]byte(data))
	meta, err := task.Meta.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}
	id, err := uuid.FromString(task.Id)
	if err != nil {
		t.Fatal(err)
	}
	if id.Version() != 1 {
		t.Fatal("expected task id to be uuid version 1")
	}
	if string(meta) != `{"id": 123}` {
		t.Fatal(`expected task meta to be {"id": 123}`)
	}
	if task.Priority != 22.5 {
		t.Fatal("expected task priority to be 22.5")
	}
	if task.Key != "tb1" {
		t.Fatal("expected task testbed to be tb1")
	}
	if !task.Created.Before(time.Now()) {
		t.Fatal("expected task created time be before now")
	}
}
