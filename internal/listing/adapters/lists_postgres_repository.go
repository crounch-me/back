package adapters

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	commonErrors "github.com/crounch-me/back/internal/common/errors"
	"github.com/crounch-me/back/internal/listing/domain/lists"
)

type ListsPostgresRepository struct {
	session *sql.DB
	schema  string
}

func NewListsPostgresRepository(session *sql.DB, schema string) (*ListsPostgresRepository, error) {
	if session == nil {
		return nil, errors.New("db session is nil")
	}

	return &ListsPostgresRepository{
		session: session,
		schema:  schema,
	}, nil
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
		return nil, commonErrors.NewError(commonErrors.UnknownErrorCode).WithCause(err)
	}

	result := make([]*lists.List, 0)
	for rows.Next() {
		if err = rows.Err(); err != nil {
			return nil, commonErrors.NewError(commonErrors.UnknownErrorCode).WithCause(err)
		}

		list := &List{}

		err = rows.Scan(&list.ID, &list.Name, &list.CreationDate, &list.ArchivationDate)
		if err != nil {
			return nil, commonErrors.NewError(commonErrors.UnknownErrorCode).WithCause(err)
		}

		domainList, err := lists.NewList(list.ID, list.Name, list.CreationDate, list.ArchivationDate)
		if err != nil {
			return nil, err
		}

		result = append(result, domainList)
	}

	return result, nil
}

func (r *ListsPostgresRepository) ReadByID(uuid string) (*lists.List, error) {
	query := fmt.Sprintf(`
    SELECT l.id, l.name, l.creation_date, l.archivation_date
    FROM %s.list l
    WHERE l.id = $1
  `, r.schema)

	row := r.session.QueryRow(query, uuid)

	list := &List{}

	err := row.Scan(&list.ID, &list.Name, &list.CreationDate, &list.ArchivationDate)
	if err != nil {
		return nil, commonErrors.NewError(commonErrors.UnknownErrorCode).WithCause(err)
	}

	return lists.NewList(list.ID, list.Name, list.CreationDate, list.ArchivationDate)
}

func (r *ListsPostgresRepository) AddContributor(listUUID, userUUID string) error {
	fmt.Println("called AddContributor")
	return nil
}

func (r *ListsPostgresRepository) GetContributorListUUIDs(userUUID string) ([]string, error) {
	fmt.Println("called GetUserListUUIDs")
	return []string{}, nil
}

func (r *ListsPostgresRepository) UpdateList(l *List) error {
	query := fmt.Sprintf(`
    UPDATE %s.list
    SET name = $1,
    archivation_date = $1
  `, r.schema)

	_, err := r.session.Exec(query, l.Name, l.ArchivationDate)

	if err != nil {
		return commonErrors.NewError(commonErrors.UnknownErrorCode).WithCause(err)
	}

	return nil
}

func (r *ListsPostgresRepository) SaveList(l *lists.List) error {
	query := fmt.Sprintf(`
    INSERT INTO %s.list(id, name, creation_date)
    VALUES ($1, $2, $3)
  `, r.schema)

	_, err := r.session.Exec(query, l.UUID(), l.Name(), l.CreationDate())
	if err != nil {
		return commonErrors.NewError(commonErrors.UnknownErrorCode).WithCause(err)
	}

	for _, c := range l.Contributors() {
		query = fmt.Sprintf(`
      INSERT INTO %s.contributor(list_id, contributor_id)
      VALUES ($1, $2)
    `, r.schema)

		_, err := r.session.Exec(query, l.UUID(), c.UUID())
		if err != nil {
			return err
		}
	}

	for _, p := range l.Products() {
		query := fmt.Sprintf(`
      INSERT INTO %s.article(article_id, list_id)
      VALUES ($1, $2)
    `, r.schema)

		_, err := r.session.Exec(query, p.UUID(), l.UUID())
		if err != nil {
			return err
		}
	}
	return nil
}
