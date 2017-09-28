package main

import (
	"flag"
	"os"
	"testing"

	arango "github.com/arangodb/go-driver"
	arangohttp "github.com/arangodb/go-driver/http"
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

func (m MockModel) FetchAll() ([]interface{}, error) {
	return make([]interface{}, 0), nil
}

func (m MockModel) Save(interface{}) (DocumentMeta, error) {
	return DocumentMeta{}, nil
}

func TestPriorityQueueModelCreate(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	model := new(PriorityQueueModel)
	if err := model.Create(); err != nil {
		t.Fatal(err)
	}
}

func TestPriorityQueueModelFetchAll(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	pq := NewPriorityQueue("test")
	model := new(PriorityQueueModel)
	if _, err := model.Save(pq); err != nil {
		t.Fatal(err)
	}
	queues, err := model.FetchAll()
	if err != nil {
		t.Fatal(err)
	}
	if queues[0].(*PriorityQueue).Key != "test" {
		t.Fatal("expected queue key to be 'test'")
	}
}

func TestPriorityQueueModelSave(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	pq := NewPriorityQueue("test")
	model := new(PriorityQueueModel)
	if _, err := model.Save(pq); err != nil {
		t.Fatal(err)
	}
	pq.Push(&Task{Id: "abc", Priority: 1.5})
	pq.Push(&Task{Id: "xyz", Priority: 2.5})
	if _, err := model.Save(pq); err != nil {
		t.Fatal(err)
	}
}
