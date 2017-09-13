package tasks_test

import (
	"flag"
	"os"
	"testing"
	"time"

	arango "github.com/arangodb/go-driver"
	arangohttp "github.com/arangodb/go-driver/http"
	"github.com/satori/go.uuid"

	"concord-pq/database"
	"concord-pq/tasks"
)

func TestMain(m *testing.M) {
	flag.Parse()
	if !testing.Short() {
		database.InitDatabase()
	}
	result := m.Run()
	if !testing.Short() {
		tearDownDatabase()
	}
	os.Exit(result)
}

func tearDownDatabase() {
	host := os.Getenv("ARANGODB_HOST")
	name := os.Getenv("ARANGODB_NAME")
	user := os.Getenv("ARANGODB_USER")
	pass := os.Getenv("ARANGODB_PASS")
	conn, err := arangohttp.NewConnection(
		arangohttp.ConnectionConfig{Endpoints: []string{host}},
	)
	if err != nil {
		panic(err)
	}
	client, err := arango.NewClient(arango.ClientConfig{
		Connection:     conn,
		Authentication: arango.BasicAuthentication(user, pass),
	})
	if err != nil {
		panic(err)
	}
	if db, err := client.Database(nil, name); err != nil {
		panic(err)
	} else {
		if err = db.Remove(nil); err != nil {
			panic(err)
		}
	}
}

type MockModel struct{}

func (m MockModel) Create() error {
	return nil
}

func (m MockModel) Save(interface{}) (arango.DocumentMeta, error) {
	return arango.DocumentMeta{}, nil
}

func TestNewTaskStat(t *testing.T) {
	stat := tasks.NewTaskStat("key", 34.5)
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

func TestTaskStatSave(t *testing.T) {
	var model database.Model
	if testing.Short() {
		model = MockModel{}
	} else {
		model = &database.TaskStatModel{}
	}
	stat := tasks.NewTaskStat("key", 34.5)
	if _, err := stat.Save(model); err != nil {
		t.Fatal(err)
	}
}

func TestNewTask(t *testing.T) {
	data := `{"meta": {"id": 123}, "Priority": 22.5, "key": "tb1"}`
	task := tasks.NewTask([]byte(data))
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
	var model database.Model
	if testing.Short() {
		model = MockModel{}
	} else {
		model = &database.TaskModel{}
	}
	task := tasks.NewTask([]byte(`{"key": "tb1", "priority": 13141.42}`))
	if _, err := task.Save(model); err != nil {
		t.Fatal(err)
	}
}

func TestTaskChangeStatus(t *testing.T) {
	var model database.Model
	if testing.Short() {
		model = MockModel{}
	} else {
		model = &database.TaskModel{}
	}
	task := tasks.NewTask([]byte(""))
	if err := task.ChangeStatus(model, tasks.StatusQueued); err != nil {
		t.Fatal(err)
	}
}
