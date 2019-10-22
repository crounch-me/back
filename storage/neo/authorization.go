package neo

import (
	"errors"

	"github.com/Sehsyha/crounch-back/model"
	"github.com/Sehsyha/crounch-back/util"

	"github.com/neo4j/neo4j-go-driver/neo4j"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

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
