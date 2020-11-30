package adapters

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/crounch-me/back/internal"
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
	query := fmt.Sprintf(`
		INSERT INTO %s."list"(id, name, creation_date)
		VALUES ($1, $2, $3)
	`, r.schema)

	_, err := r.session.Exec(query, list.UUID(), list.Name(), list.CreationDate())

	if err != nil {
		return internal.NewError(internal.UnknownErrorCode).WithCause(err)
	}

	return nil
}

func (r *ListsPostgresRepository) ReadByIDs(uuids []string) ([]*lists.List, error) {
	query := fmt.Sprintf(`
    SELECT l.id, l.name, l.creation_date, l.archivation_date
    FROM %s.list l
    WHERE l.id IN ($1)
  `, r.schema)

	listUUIDs := strings.Join(uuids, ",")

	rows, err := r.session.Query(query, listUUIDs)
	defer rows.Close()
	if err != nil {
		return nil, internal.NewError(internal.UnknownErrorCode).WithCause(err)
	}

	result := make([]*lists.List, 0)
	for rows.Next() {
		if err = rows.Err(); err != nil {
			return nil, internal.NewError(internal.UnknownErrorCode).WithCause(err)
		}

		list := &List{}

		err = rows.Scan(&list.ID, &list.Name, &list.CreationDate, &list.ArchivationDate)
		if err != nil {
			return nil, internal.NewError(internal.UnknownErrorCode).WithCause(err)
		}

		domainList, err := lists.NewList(list.ID, list.Name, list.CreationDate, list.ArchivationDate)
		if err != nil {
			return nil, err
		}

		result = append(result, domainList)
	}

	return result, nil
}
