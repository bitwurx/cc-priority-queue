package database_test

import (
	"flag"
	"os"
	"testing"

	arango "github.com/arangodb/go-driver"
	arangohttp "github.com/arangodb/go-driver/http"

	"concord-pq/database"
	"concord-pq/tasks"
)

func TestMain(m *testing.M) {
	flag.Parse()
	if !testing.Short() {
		database.InitDatabase()
		result := m.Run()
		tearDownDatabase()
		os.Exit(result)
	}
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

func TestTaskStatModelCreate(t *testing.T) {
	model := new(database.TaskStatModel)
	if err := model.Create(); err != nil {
		t.Fatal(err)
	}
}

func TestTaskStatModelSave(t *testing.T) {
	taskStat := tasks.NewTaskStat("key", 13.5)
	model := new(database.TaskStatModel)
	if _, err := model.Save(taskStat); err != nil {
		t.Fatal(err)
	}
}

func TestTaskModelCreate(t *testing.T) {
	model := new(database.TaskModel)
	if err := model.Create(); err != nil {
		t.Fatal(err)
	}
}

func TestTaskModelSave(t *testing.T) {
	task := tasks.NewTask([]byte(`{"meta": {"id": 123}, "Priority": 22.5, "key": "tb1"}`))
	model := new(database.TaskModel)
	if _, err := model.Save(task); err != nil {
		t.Fatal(err)
	}
}
