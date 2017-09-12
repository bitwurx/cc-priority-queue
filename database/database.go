package database

import (
	"context"
	"log"
	"os"
	"time"

	arango "github.com/arangodb/go-driver"
	arangohttp "github.com/arangodb/go-driver/http"
)

var db arango.Database
var models = make([]Model, 0)

type Model interface {
	CreateCollection(Database)
}

type Collection interface {
	CreateDocument(context.Context, interface{}) (arango.DocumentMeta, error)
	ReadDocument(context.Context, string, interface{}) (arango.DocumentMeta, error)
	UpdateDocument(context.Context, string, interface{}) (arango.DocumentMeta, error)
}

type Database interface {
	Collection(context.Context, string) (arango.Collection, error)
	CreateCollection(context.Context, string, *arango.CreateCollectionOptions) (arango.Collection, error)
}

func InitDatabase() {
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

	for _, model := range models {
		model.CreateCollection(db)
	}
}

func AddModel(model Model) {
	models = append(models, model)
}

func GetCollection(name string) (arango.Collection, error) {
	col, err := db.Collection(nil, name)
	if arango.IsNotFound(err) {
		col, err = db.CreateCollection(nil, name, nil)
	}
	return col, err
}
