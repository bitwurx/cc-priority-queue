package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	arango "github.com/arangodb/go-driver"
	arangohttp "github.com/arangodb/go-driver/http"
)

const (
	CollectionPriorityQueues = "priority_queues" // the name of the priority queues database collection.
)

var db arango.Database // package local arango database instance.

// DocumentMeta contains meta data for an arango document
type DocumentMeta struct {
	Id arango.DocumentID
}

// Model contains methods for interacting with database collections.
type Model interface {
	Create() error
	FetchAll() ([]interface{}, error)
	Save(interface{}) (DocumentMeta, error)
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

// FetchAll gets all documents from the priority queues collection.
func (model *PriorityQueueModel) FetchAll() ([]interface{}, error) {
	queues := make([]interface{}, 0)
	query := fmt.Sprintf("FOR q IN %s RETURN q", CollectionPriorityQueues)
	cursor, err := db.Query(nil, query, nil)
	if err != nil {
		return nil, err
	}
	defer cursor.Close()
	for {
		q := new(PriorityQueue)
		_, err := cursor.ReadDocument(nil, q)
		if arango.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			return nil, err
		}
		queues = append(queues, q)
	}
	return queues, nil
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
	}
	for _, model := range models {
		if err := model.Create(); err != nil {
			panic(err)
		}
	}
}
