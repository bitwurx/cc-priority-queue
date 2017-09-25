package main

import (
	"testing"
	"time"

	"github.com/satori/go.uuid"
)

type MockModel struct{}

func (m MockModel) Create() error {
	return nil
}

func (m MockModel) Query(q string, vars interface{}) ([]interface{}, error) {
	return make([]interface{}, 0), nil
}

func (m MockModel) Save(interface{}) (DocumentMeta, error) {
	return DocumentMeta{}, nil
}

func TestNewTaskStat(t *testing.T) {
	stat := NewTaskStat("key", 34.5)
	if !stat.Created.Before(time.Now()) {
		t.Fatal("expected task stat created timestamp to be before now")
	}
	if stat.RunTime != 34.5 {
		t.Fatal("expected task stat runtime to be 34.5")
	}
	if stat.Key != "key" {
		t.Fatal("expected task stat key to be 'key'")
	}
}

func TestTaskStatSave(t *testing.T) {
	var model Model
	if testing.Short() {
		model = MockModel{}
	} else {
		model = &TaskStatModel{}
	}
	stat := NewTaskStat("key", 34.5)
	if _, err := stat.Save(model); err != nil {
		t.Fatal(err)
	}
}

func TestNewTask(t *testing.T) {
	data := `{"meta": {"id": 123}, "Priority": 22.5, "key": "tb1"}`
	task := NewTask([]byte(data))
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

func TestTaskSave(t *testing.T) {
	var model Model
	if testing.Short() {
		model = MockModel{}
	} else {
		model = &TaskModel{}
	}
	task := NewTask([]byte(`{"key": "tb1", "priority": 13141.42}`))
	if _, err := task.Save(model); err != nil {
		t.Fatal(err)
	}
}

func TestTaskChangeStatus(t *testing.T) {
	var model Model
	if testing.Short() {
		model = MockModel{}
	} else {
		model = &TaskModel{}
	}
	task := NewTask([]byte(""))
	if err := task.ChangeStatus(model, StatusQueued); err != nil {
		t.Fatal(err)
	}
}

type MockTaskStatGetRunTimeModel struct {
	MockModel
	Documents map[string][]interface{}
}

func (m MockTaskStatGetRunTimeModel) Save(t interface{}) (DocumentMeta, error) {
	v, _ := t.(*TaskStat)
	m.Documents[v.Key] = append(m.Documents[v.Key], v)
	return DocumentMeta{}, nil
}

func (m MockTaskStatGetRunTimeModel) Query(q string, vars interface{}) ([]interface{}, error) {
	v, _ := vars.(map[string]interface{})
	return m.Documents[v["key"].(string)][5:15], nil
}

func TestTaskGetRunTime(t *testing.T) {
	var taskModel Model
	var taskStatModel Model
	if testing.Short() {
		taskModel = MockModel{}
		taskStatModel = MockTaskStatGetRunTimeModel{MockModel{}, make(map[string][]interface{})}
	} else {
		taskModel = &TaskModel{}
		taskStatModel = &TaskStatModel{}
	}
	task := NewTask([]byte(`{"key": "my-task"}`))
	if _, err := task.Save(taskModel); err != nil {
		t.Fatal(err)
	}
	stats := []*TaskStat{
		NewTaskStat("my-task", 1.0),
		NewTaskStat("my-task", 1.0),
		NewTaskStat("my-task", 1.0),
		NewTaskStat("my-task", 1.0),
		NewTaskStat("my-task", 1.0),
		NewTaskStat("my-task", 1.0),
		NewTaskStat("my-task", 1.0),
		NewTaskStat("my-task", 5.0),
		NewTaskStat("my-task", 1.0),
		NewTaskStat("my-task", 3.0),
		NewTaskStat("my-task", 1.0),
		NewTaskStat("my-task", 7.0),
		NewTaskStat("my-task", 1.0),
		NewTaskStat("my-task", 2.0),
		NewTaskStat("my-task", 1.0),
	}
	for _, stat := range stats {
		if _, err := stat.Save(taskStatModel); err != nil {
			t.Fatal(err)
		}
	}
	if runTime, err := task.GetAverageRunTime(taskStatModel); err != nil {
		t.Fatal(err)
	} else {
		if runTime != 2.3 {
			t.Fatal("expected average runtime to be 2.3")
		}
	}
}
