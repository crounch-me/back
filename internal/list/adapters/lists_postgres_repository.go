package adapters

import "database/sql"

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
