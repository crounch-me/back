package adapters

import (
	"database/sql"
	"errors"
	"fmt"
)

type ContributorsPostgresRepository struct {
	session *sql.DB
	schema  string
}

func NewContributorsPostgresRepository(session *sql.DB, schema string) (*ContributorsPostgresRepository, error) {
	if session == nil {
		return nil, errors.New("db session is nil")
	}

	return &ContributorsPostgresRepository{
		session: session,
		schema:  schema,
	}, nil
}

func (r ContributorsPostgresRepository) AddContributor(listUUID, userUUID string) error {
	fmt.Println("called AddContributor")
	return nil
}

func (r ContributorsPostgresRepository) GetUserListUUIDs(userUUID string) ([]string, error) {
	fmt.Println("called GetUserListUUIDs")
	return []string{}, nil
}
