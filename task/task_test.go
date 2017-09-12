package task_test

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"os"
	"testing"
	"time"

	arango "github.com/arangodb/go-driver"
	arangohttp "github.com/arangodb/go-driver/http"
	"github.com/satori/go.uuid"

	. "concord-pq/database"
	. "concord-pq/task"
)

func TestMain(m *testing.M) {
	flag.Parse()
	if !testing.Short() {
		InitDatabase()
	}
	result := m.Run()
	if !testing.Short() {
		tearDownDatabase()
	}
	os.Exit(result)
}

// tearDownDatabase drops the test__concord_pq database
func tearDownDatabase() {
	host := os.Getenv("ARANGODB_HOST")
	name := os.Getenv("ARANGODB_NAME")
	user := os.Getenv("ARANGODB_USER")
	pass := os.Getenv("ARANGODB_PASS")
	conn, err := arangohttp.NewConnection(
		arangohttp.ConnectionConfig{Endpoints: []string{host}},
	)
	if err != nil {
		log.Fatal(err)
	}
	client, err := arango.NewClient(arango.ClientConfig{
		Connection:     conn,
		Authentication: arango.BasicAuthentication(user, pass),
	})
	if err != nil {
		log.Fatal(err)
	}
	if db, err := client.Database(nil, name); err != nil {
		log.Fatal(err)
	} else {
		if err = db.Remove(nil); err != nil {
			log.Fatal(err)
		}
	}
}

type MockDatabase struct {
	Collections map[string]Collection
}

func (db *MockDatabase) Collection(ctx context.Context, name string) (arango.Collection, error) {
	return nil, nil
}

func (db *MockDatabase) CreateCollection(ctx context.Context, name string, opts *arango.CreateCollectionOptions) (arango.Collection, error) {
	if _, ok := db.Collections[name]; ok {
		return nil, arango.ArangoError{Code: 409}
	}
	db.Collections[name] = nil
	return nil, nil
}

type MockTaskStatsCollection struct {
	Documents map[string][]byte
}

func (col *MockTaskStatsCollection) CreateDocument(ctx context.Context, taskStat interface{}) (arango.DocumentMeta, error) {
	v, _ := taskStat.(*TaskStat)
	if _, ok := col.Documents[v.Key]; ok {
		return arango.DocumentMeta{}, arango.ArangoError{Code: 409}
	}
	data, err := json.Marshal(taskStat)
	if err != nil {
		return arango.DocumentMeta{}, err
	}
	col.Documents[v.Key] = data
	return arango.DocumentMeta{}, nil
}

func (col *MockTaskStatsCollection) ReadDocument(ctx context.Context, key string, taskStat interface{}) (arango.DocumentMeta, error) {
	doc, ok := col.Documents[key]
	if !ok {
		return arango.DocumentMeta{}, arango.ArangoError{Code: 404}
	}
	json.Unmarshal(doc, taskStat)
	return arango.DocumentMeta{}, nil
}

func (col *MockTaskStatsCollection) UpdateDocument(ctx context.Context, key string, patch interface{}) (arango.DocumentMeta, error) {
	return arango.DocumentMeta{}, nil
}

func TestNewTaskStat(t *testing.T) {
	stat := NewTaskStat("key", 34.5)
	if !stat.Created.Before(time.Now()) {
		t.Fatal("expected task stat created timestamp to be before now")
	}
	if stat.Runtime != 34.5 {
		t.Fatal("expected task stat runtime to be 34.5")
	}
	if stat.Key != "key" {
		t.Fatal("expected task stat key to be 'key'")
	}
}

func TestTaskStatCreateCollection(t *testing.T) {
	db := &MockDatabase{make(map[string]Collection)}
	stat := NewTaskStat("key", 0)
	stat.CreateCollection(db)
}

func TestTaskStatSave(t *testing.T) {
	var taskStats Collection
	if testing.Short() {
		taskStats = &MockTaskStatsCollection{make(map[string][]byte)}
	} else {
		var err error
		taskStats, err = GetCollection("task_stats")
		if err != nil {
			t.Fatal(err)
		}
	}
	stat := NewTaskStat("key", 34.5)
	stat.Save(taskStats)
}

type MockTasksCollection struct {
	Documents map[string][]byte
}

func (col *MockTasksCollection) CreateDocument(ctx context.Context, task interface{}) (arango.DocumentMeta, error) {
	v, _ := task.(*Task)
	if _, ok := col.Documents[v.Id]; ok {
		return arango.DocumentMeta{}, arango.ArangoError{Code: 409}
	}
	data, err := json.Marshal(task)
	if err != nil {
		return arango.DocumentMeta{}, err
	}
	col.Documents[v.Id] = data
	return arango.DocumentMeta{}, nil
}

func (col *MockTasksCollection) ReadDocument(ctx context.Context, key string, task interface{}) (arango.DocumentMeta, error) {
	doc, ok := col.Documents[key]
	if !ok {
		return arango.DocumentMeta{}, arango.ArangoError{Code: 404}
	}
	json.Unmarshal(doc, task)
	return arango.DocumentMeta{}, nil
}

func (col *MockTasksCollection) UpdateDocument(ctx context.Context, key string, patch interface{}) (arango.DocumentMeta, error) {
	v, _ := patch.(map[string]interface{})
	doc, _ := col.Documents[key]
	task := new(Task)
	json.Unmarshal(doc, task)
	task.Status = v["Status"].(int)
	data, err := json.Marshal(task)
	if err != nil {
		return arango.DocumentMeta{}, err
	}
	doc = data
	return arango.DocumentMeta{}, nil
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

func TestTaskCreateCollection(t *testing.T) {
	db := &MockDatabase{make(map[string]Collection)}
	task := NewTask([]byte(`{"key": "xyz", "priority": 32.5}`))
	task.CreateCollection(db)
}

func TestTaskSave(t *testing.T) {
	var tasks Collection
	if testing.Short() {
		tasks = &MockTasksCollection{make(map[string][]byte)}
	} else {
		var err error
		tasks, err = GetCollection("tasks")
		if err != nil {
			t.Fatal(err)
		}
	}
	task := NewTask([]byte(`{"key": "tb1", "priority": 13141.42}`))
	task.Save(tasks)
	if err := task.Save(tasks); err != nil {
		t.Fatal(err)
	}
	readTask := new(Task)
	if _, err := tasks.ReadDocument(nil, task.Id, readTask); err != nil {
		t.Fatal(err)
	}
	if readTask.Priority != 13141.42 {
		t.Fatal("expected task priority to be 13141.42")
	}
	if readTask.Status != StatusPending {
		t.Fatal("expected task status to be Pending")
	}
}

func TestTaskChangeStatus(t *testing.T) {
	var tasks Collection
	if testing.Short() {
		tasks = &MockTasksCollection{make(map[string][]byte)}
	} else {
		var err error
		tasks, err = GetCollection("tasks")
		if err != nil {
			t.Fatal(err)
		}
	}
	task := NewTask([]byte(""))
	if err := task.ChangeStatus(tasks, StatusQueued); err != nil {
		t.Fatal(err)
	}
	readTask := new(Task)
	if _, err := tasks.ReadDocument(nil, task.Id, readTask); err != nil {
		t.Fatal(err)
	}
	if readTask.Status != StatusQueued {
		t.Fatal("expected task status to be queued")
	}
}
