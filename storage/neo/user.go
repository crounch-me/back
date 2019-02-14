package neo

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"

	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"

	"github.com/Sehsyha/crounch-back/model"
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

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
