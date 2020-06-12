package postgres

import (
	"database/sql"
	"fmt"

	"github.com/crounch-me/back/model"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

// CreateUser inserts a new user with hashed password
func (s *PostgresStorage) CreateUser(user *model.User) error {
	log.WithField("email", user.Email).Debug("Creating user")

	user.ID = uuid.NewV4().String()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	stringPassword := string(hashedPassword)

	user.Password = &stringPassword

	query := fmt.Sprintf(`
		INSERT INTO %s."user"(id, email, password)
		VALUES ($1, $2, $3)
	`, s.schema)

	_, err = s.session.Exec(query, user.ID, user.Email, user.Password)

	if err != nil {
		log.WithError(err).Error("Unable to create user")
		return err
	}

	return nil
}

// GetUserByEmail find the user with his email
func (s *PostgresStorage) GetUserByEmail(email string) (*model.User, error) {
	log.WithField("email", email).Debug("Getting user by email")

	query := fmt.Sprintf(`
		SELECT id, password
		FROM %s."user"
		WHERE LOWER("user".email) = LOWER($1)
	`, s.schema)

	row := s.session.QueryRow(query, email)

	user := &model.User{}

	err := row.Scan(&user.ID, &user.Password)

	if err == sql.ErrNoRows {
		return nil, model.NewDatabaseError(model.ErrNotFound, nil)
	}

	if err != nil {
		log.WithError(err).Error("Unable to find user by email")
		return nil, err
	}

	return user, nil
}

// GetUserIDByToken find the user with his token
func (s *PostgresStorage) GetUserIDByToken(token string) (*string, error) {
	log.WithField("token", token).Debug("Getting user by token")

	query := fmt.Sprintf(`
		SELECT id
    FROM %s."user"
    LEFT JOIN %s."authorization" ON "authorization".user_id = "user".id
		WHERE "authorization".token = $1
	`, s.schema, s.schema)

	row := s.session.QueryRow(query, token)

	var id *string

	err := row.Scan(&id)

	if err == sql.ErrNoRows {
		return nil, model.NewDatabaseError(model.ErrNotFound, nil)
	}

	if err != nil {
		log.WithError(err).Error("Unable to find user by token")
		return nil, err
	}

	return id, nil
}
