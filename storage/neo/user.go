package neo

import (
	"errors"
	"fmt"
	"math/rand"

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
	log.WithField("email", u.Email).Info("Creating user")
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

func (ns *NeoStorage) CreateAuthorization(u *model.User) (*model.Authorization, error) {
	log.WithField("email", u.Email).Info("Creating authorization for user")

	session, err := ns.driver.Session(neo4j.AccessModeWrite)
	if err != nil {
		return nil, err
	}

	query := `
		MATCH (u: USER {email: $email})
		RETURN u.password, u.id
	`

	result, err := session.Run(query, map[string]interface{}{
		"email": u.Email,
	})
	if err != nil {
		return nil, err
	}

	isResult := result.Next()
	if err := result.Err(); err != nil {
		return nil, err
	}

	if !isResult {
		return nil, errors.New("user not found")
	}

	hashedPassword := result.Record().GetByIndex(0).(string)
	u.ID = result.Record().GetByIndex(1).(string)
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(u.Password))
	if err != nil {
		return nil, errors.New("wrong password")
	}

	query = `
		MATCH (u: USER {id: $id})
		CREATE (a: AUTHORIZATION {token: $token})
		CREATE (u)-[:IS_AUTHORIZED_BY]->(a)
		RETURN ID(a)
	`

	token := generateToken(42)

	result, err = session.Run(query, map[string]interface{}{
		"id":    u.ID,
		"token": token,
	})

	isResult = result.Next()
	if err := result.Err(); err != nil {
		return nil, err
	}

	if !isResult {
		return nil, errors.New("error occured while creating authorization")
	}

	if err != nil {
		return nil, err
	}

	authorization := &model.Authorization{
		AccessToken: token,
	}

	return authorization, nil
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func generateToken(size int) string {
	b := make([]byte, size)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
