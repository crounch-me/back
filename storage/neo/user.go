package neo

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	uuid "github.com/satori/go.uuid"

	"github.com/Sehsyha/crounch-back/model"
	"github.com/Sehsyha/crounch-back/storage"
	"github.com/Sehsyha/crounch-back/util"
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

func (ns *NeoStorage) CreateUser(u *model.User) error {
	log.WithField("email", u.Email).Debug("Creating user")
	session, err := ns.driver.Session(neo4j.AccessModeWrite)
	if err != nil {
		return err
	}
	defer session.Close()

	u.ID = uuid.NewV4().String()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = string(hashedPassword)

	query := `
		CREATE (u: USER {email: $email, id: $id, password: $password})
		RETURN u.id
	`

	result, err := session.Run(query, map[string]interface{}{
		"email":    u.Email,
		"id":       u.ID,
		"password": u.Password,
	})
	if err != nil {
		return err
	}

	isResult := result.Next()
	if err = result.Err(); err != nil {
		return err
	}

	if !isResult {
		return model.NewDatabaseError(model.ErrCreation, nil)
	}

	for result.Next() {
		fmt.Printf("Created %s \n", result.Record().GetByIndex(0))
	}

	return nil
}

func (ns *NeoStorage) GetUserByEmail(email string) (*model.User, error) {
	log.WithField("email", email).Debug("Getting user by email")

	session, err := ns.driver.Session(neo4j.AccessModeRead)
	if err != nil {
		return nil, err
	}

	query := `
		MATCH (u: USER {email: $email})
		RETURN u.id, u.password
	`

	result, err := session.Run(query, map[string]interface{}{
		"email": email,
	})
	if err != nil {
		return nil, err
	}

	isResult := result.Next()
	if err := result.Err(); err != nil {
		return nil, err
	}

	if !isResult {
		return nil, model.NewDatabaseError(model.ErrNotFound, nil)
	}

	u := &model.User{
		ID:       result.Record().GetByIndex(0).(string),
		Email:    email,
		Password: result.Record().GetByIndex(1).(string),
	}

	return u, nil
}

func (ns *NeoStorage) CreateAuthorization(u *model.User) (*model.Authorization, error) {
	log.WithField("email", u.Email).Debug("Creating authorization for user")

	session, err := ns.driver.Session(neo4j.AccessModeWrite)
	if err != nil {
		return nil, err
	}
	defer session.Close()

	foundUser, err := ns.GetUserByEmail(u.Email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(u.Password))
	if err != nil {
		return nil, model.NewDatabaseError(model.ErrWrongPassword, nil)
	}

	query := `
		MATCH (u: USER {id: $id})
		CREATE (a: AUTHORIZATION {token: $token})
		CREATE (u)-[:IS_AUTHORIZED_BY]->(a)
		RETURN ID(a)
	`

	token := util.GenerateToken()

	result, err := session.Run(query, map[string]interface{}{
		"id":    foundUser.ID,
		"token": token,
	})

	isResult := result.Next()
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
