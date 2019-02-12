package neo

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"

	uuid "github.com/satori/go.uuid"

	"github.com/Sehsyha/crounch-back/model"
	"github.com/Sehsyha/crounch-back/storage"
	"github.com/neo4j/neo4j-go-driver/neo4j"
	log "github.com/sirupsen/logrus"
)

type NeoStorage struct {
	driver neo4j.Driver
}

func NewNeoStorage() storage.Storage {
	driver, err := neo4j.NewDriver("bolt://localhost:7687", neo4j.BasicAuth("test", "test", ""))
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

func (ns *NeoStorage) CreateUser(u *model.User) error {
	log.Info("Creating user")
	session, err := ns.driver.Session(neo4j.AccessModeWrite)
	if err != nil {
		return err
	}

	u.ID = uuid.NewV4().String()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = string(hashedPassword)

	query := `
		CREATE (u: USER {email: $email, id: $id, password: $password})
	`

	result, err := session.Run(query, map[string]interface{}{
		"email":    u.Email,
		"id":       u.ID,
		"password": u.Password,
	})
	if err != nil {
		return err
	}

	for result.Next() {
		fmt.Printf("Created %s \n", result.Record().GetByIndex(0))
	}

	return nil
}
