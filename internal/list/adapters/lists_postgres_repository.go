package adapters

import (
	"database/sql"
	"fmt"

	"github.com/crounch-me/back/internal/list/domain/lists"
)

type ListsPostgresRepository struct {
	session *sql.DB
	schema  string
}

func NewListsPostgresRepository(session *sql.DB, schema string) *ListsPostgresRepository {
	return &ListsPostgresRepository{
		session: session,
		schema:  schema,
	}
}

func (r *ListsPostgresRepository) AddList(list *lists.List) error {
	fmt.Println("called AddList")
	return nil
}

func (r *ListsPostgresRepository) ReadByIDs(uuids []string) ([]*lists.List, error) {
	fmt.Println("called ReadByIDs")
	return nil, nil
}
