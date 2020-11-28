package postgres

import (
	"fmt"

	"github.com/crounch-me/back/internal"
)

func (s *PostgresStorage) CreateContributor(listID, userID string) *internal.Error {
	query := fmt.Sprintf(`
    INSERT INTO %s."contributor" (list_id, user_id)
    VALUES ($1, $2)
  `, s.schema)

	_, err := s.session.Exec(query, listID, userID)

	if err != nil {
		return internal.NewError(internal.UnknownErrorCode).WithCause(err)
	}

	return nil
}

func (s *PostgresStorage) GetContributorsIDs(listID string) ([]string, *internal.Error) {
	query := fmt.Sprintf(`
    SELECT c.user_id
    FROM %s.contributor c
    WHERE c.list_id = $1
  `, s.schema)

	rows, err := s.session.Query(query, listID)
	defer rows.Close()

	if err != nil {
		return nil, internal.NewError(internal.UnknownErrorCode).WithCause(err)
	}

	contributorIDs := make([]string, 0)
	for rows.Next() {
		if err = rows.Err(); err != nil {
			return nil, internal.NewError(internal.UnknownErrorCode).WithCause(err)
		}

		contributorID := ""

		err = rows.Scan(&contributorID)
		if err != nil {
			return nil, internal.NewError(internal.UnknownErrorCode).WithCause(err)
		}

		contributorIDs = append(contributorIDs, contributorID)
	}

	return contributorIDs, nil
}
