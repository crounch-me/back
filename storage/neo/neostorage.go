package neo

import (
	"github.com/Sehsyha/crounch-back/storage"
	"github.com/neo4j/neo4j-go-driver/neo4j"
	log "github.com/sirupsen/logrus"
)

type NeoStorage struct {
	driver neo4j.Driver
}

func NewNeoStorage() storage.Storage {
	driver, err := neo4j.NewDriver("bolt://database:7687", neo4j.BasicAuth("test", "test", ""))
	if err != nil {
		log.Fatal(err)
	}
	session, err := driver.Session(neo4j.AccessModeRead)
	if err != nil {
		log.Fatal(err)
	}
	_, err = session.Run("MATCH (n:TEST)", map[string]interface{}{})
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()
	return &NeoStorage{
		driver: driver,
	}
}
