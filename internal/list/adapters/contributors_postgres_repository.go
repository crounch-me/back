package adapters

import (
	"database/sql"
	"fmt"
)

type ContributorsPostgresRepository struct {
	session *sql.DB
	schema  string
}

func NewContributorsPostgresRepository(session *sql.DB, schema string) *ContributorsPostgresRepository {
	return &ContributorsPostgresRepository{
		session: session,
		schema:  schema,
	}
}

func (r ContributorsPostgresRepository) AddContributor(listUUID, userUUID string) error {
	fmt.Println("called AddContributor")
	return nil
}

func (r ContributorsPostgresRepository) GetUserListUUIDs(userUUID string) ([]string, error) {
	fmt.Println("called GetUserListUUIDs")
	return []string{}, nil
}
