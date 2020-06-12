package postgres

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"

	"github.com/crounch-me/back/model"
	"github.com/crounch-me/back/util"
)

// CreateAuthorization creates a user id and token couple
func (s *PostgresStorage) CreateAuthorization(user *model.User) (*model.Authorization, error) {
	log.WithField("email", user.Email).Debug("Creating authorization for user")

	foundUser, err := s.GetUserByEmail(user.Email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(*foundUser.Password), []byte(*user.Password))
	if err != nil {
		return nil, model.NewDatabaseError(model.ErrWrongPassword, nil)
	}

	query := fmt.Sprintf(`
		INSERT INTO %s."authorization" (user_id, token)
		VALUES ($1, $2)
	`, s.schema)

	token := util.GenerateToken()
	_, err = s.session.Exec(query, foundUser.ID, token)

	if err != nil {
		log.WithError(err).Error("Unable to create authorization")
		return nil, err
	}

	authorization := &model.Authorization{
		AccessToken: token,
	}

	return authorization, nil
}
