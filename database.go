package main

import (
	"encoding/json"
	"os"
	"time"

	arango "github.com/arangodb/go-driver"
	arangohttp "github.com/arangodb/go-driver/http"
)

const (
	CollectionPriorityQueues = "priority_queues" // the name of the priority queues database collection.
	CollectionTasks          = "tasks"           // the name of the tasks database collection.
	CollectionTaskStats      = "task_stats"      // the name of the task stats database collection.
)

var db arango.Database // package local arango database instance.

// DocumentMeta contains meta data for an arango document
type DocumentMeta struct {
	Id arango.DocumentID
}

// Model contains methods for interacting with database collections.
type Model interface {
	Create() error
	Query(string, interface{}) ([]interface{}, error)
	Save(interface{}) (DocumentMeta, error)
}

// TaskStatModel represents a task stat collection model.
type TaskStatModel struct{}

// Create creates the task_stats collection and creates a persistent index on
// the Created field in the arangodb database.
func (model *TaskStatModel) Create() error {
	col, err := db.CreateCollection(nil, CollectionTaskStats, nil)
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

// Query runs the AQL query against the task stat model collection.
func (model *TaskStatModel) Query(q string, vars interface{}) ([]interface{}, error) {
	taskStats := make([]interface{}, 0)
	cursor, err := db.Query(nil, q, vars.(map[string]interface{}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close()
	for {
		taskStat := new(TaskStat)
		_, err := cursor.ReadDocument(nil, taskStat)
		if arango.IsNoMoreDocuments(err) {
			break
		}
		if err != nil {
			return nil, err
		}
		taskStats = append(taskStats, taskStat)
	}
	return taskStats, nil
}

// Save creates a document in the task stats collection.
func (model *TaskStatModel) Save(taskStat interface{}) (DocumentMeta, error) {
	col, err := db.Collection(nil, CollectionTaskStats)
	if err != nil {
		return DocumentMeta{}, err
	}
	meta, err := col.CreateDocument(nil, taskStat)
	if err != nil {
		return DocumentMeta{}, err
	}
	return DocumentMeta{Id: meta.ID}, nil
}

// TaskModel represents a task collection model.
type TaskModel struct{}

// Create creates the tasks collection in the arangodb database.
func (model *TaskModel) Create() error {
	_, err := db.CreateCollection(nil, CollectionTasks, nil)
	if err != nil && arango.IsConflict(err) {
		return nil
	}
	return err
}

// Query runs the AQL query against the task model collection.
func (model *TaskModel) Query(q string, vars interface{}) ([]interface{}, error) {
	return make([]interface{}, 0), nil
}

// Save creates a document in the tasks collection.
func (model *TaskModel) Save(task interface{}) (DocumentMeta, error) {
	var meta arango.DocumentMeta
	col, err := db.Collection(nil, CollectionTasks)
	if err != nil {
		return DocumentMeta{}, err
	}
	meta, err = col.CreateDocument(nil, task)
	if arango.IsConflict(err) {
		v, _ := task.(*Task)
		patch := map[string]interface{}{"status": v.Status}
		meta, err = col.UpdateDocument(nil, v.Id, patch)
		if err != nil {
			return DocumentMeta{}, err
		}
	} else if err != nil {
		return DocumentMeta{}, err
	}
	return DocumentMeta{Id: meta.ID}, nil
}

// PriorityQueueModel represents a priority queue collection model.
type PriorityQueueModel struct{}

// Create creates the priority queues collection in the arangodb database.
func (model *PriorityQueueModel) Create() error {
	_, err := db.CreateCollection(nil, CollectionPriorityQueues, nil)
	if err != nil && arango.IsConflict(err) {
		return nil
	}
	return err
}

// Query runs the AQL query against the priority queues model collection.
func (model *PriorityQueueModel) Query(q string, vars interface{}) ([]interface{}, error) {
	return make([]interface{}, 0), nil
}

// Query runs the AQL query against the priority queues model collection.
func (model *PriorityQueueModel) Save(pq interface{}) (DocumentMeta, error) {
	var meta arango.DocumentMeta
	var doc struct {
		Key   string      `json:"_key"`
		Count int         `json:"count"`
		Heap  interface{} `json:"heap"`
	}
	col, err := db.Collection(nil, CollectionPriorityQueues)
	if err != nil {
		return DocumentMeta{}, err
	}
	data, err := json.Marshal(pq.(*PriorityQueue))
	if err != nil {
		return DocumentMeta{}, err
	}
	if err := json.Unmarshal(data, &doc); err != nil {
		return DocumentMeta{}, err
	}
	meta, err = col.CreateDocument(nil, doc)
	if arango.IsConflict(err) {
		patch := map[string]interface{}{
			"count": doc.Count,
			"heap":  doc.Heap,
		}
		meta, err = col.UpdateDocument(nil, doc.Key, patch)
		if err != nil {
			return DocumentMeta{}, err
		}
	} else if err != nil {
		return DocumentMeta{}, err
	}
	return DocumentMeta{Id: meta.ID}, nil
	return DocumentMeta{}, nil
}

// InitDatabase connects to the arangodb and creates the collections from the
// provided models.
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
		&PriorityQueueModel{},
		&TaskModel{},
		&TaskStatModel{},
	}
	for _, model := range models {
		if err := model.Create(); err != nil {
			panic(err)
		}
	}
}
