package database

import (
	"os"
	"time"

	arango "github.com/arangodb/go-driver"
	arangohttp "github.com/arangodb/go-driver/http"
)

var db arango.Database // package local arango database instance

// Model contains methods for interacting with database collections
type Model interface {
	Create() error
	Save(interface{}) (arango.DocumentMeta, error)
}

// TaskStatModel represents a task stat collection model
type TaskStatModel struct{}

// Create creates the task_stats collection and creates a persistent index on
// the Created field in the arangodb database
func (model *TaskStatModel) Create() error {
	col, err := db.CreateCollection(nil, "task_stats", nil)
	if err != nil {
		if arango.IsConflict(err) {
			return nil
		}
		return err
	}
	_, _, err = col.EnsurePersistentIndex(nil, []string{"Created"}, nil)
	if err != nil {
		return err
	}
	return err
}

// Save creates a document in the task stats collection
func (model *TaskStatModel) Save(taskStat interface{}) (arango.DocumentMeta, error) {
	col, err := db.Collection(nil, "task_stats")
	if err != nil {
		return arango.DocumentMeta{}, err
	}
	return col.CreateDocument(nil, taskStat)
}

// TaskModel represents a task collection model
type TaskModel struct{}

// Create creates the tasks collection in the arangodb database
func (model *TaskModel) Create() error {
	_, err := db.CreateCollection(nil, "tasks", nil)
	if err != nil && arango.IsConflict(err) {
		return nil
	}
	return err
}

// Save creates a document in the tasks collection
func (model *TaskModel) Save(task interface{}) (arango.DocumentMeta, error) {
	col, err := db.Collection(nil, "tasks")
	if err != nil {
		return arango.DocumentMeta{}, err
	}
	return col.CreateDocument(nil, task)
}

// InitDatabase connects to the arangodb and creates the collections from the
// provided models
func InitDatabase() {
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

	for {
		if exists, err := client.DatabaseExists(nil, name); err == nil {
			if !exists {
				db, err = client.CreateDatabase(nil, name, nil)
			} else {
				db, err = client.Database(nil, name)
			}
			if err == nil {
				break
			}
		}
		time.Sleep(time.Second * 1)
	}

	models := []Model{
		&TaskModel{},
		&TaskStatModel{},
	}
	for _, model := range models {
		if err := model.Create(); err != nil {
			panic(err)
		}
	}
}
