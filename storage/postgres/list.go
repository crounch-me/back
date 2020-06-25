package postgres

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/crounch-me/back/domain"
	"github.com/crounch-me/back/domain/lists"
	"github.com/crounch-me/back/domain/users"
)

// CreateList inserts a new list
func (s *PostgresStorage) CreateList(id, name, ownerID string, creationDate time.Time) *domain.Error {
	query := fmt.Sprintf(`
		INSERT INTO %s."list"(id, name, user_id, creation_date)
		VALUES ($1, $2, $3, $4)
	`, s.schema)

	_, err := s.session.Exec(query, id, name, ownerID, creationDate)

	if err != nil {
		return domain.NewErrorWithCause(domain.UnknownErrorCode, err)
	}

	return nil
}

// GetOwnerLists get all owner's lists
func (s *PostgresStorage) GetOwnersLists(ownerID string) ([]*lists.List, *domain.Error) {
	query := fmt.Sprintf(`
    SELECT l.id, l.name, l.creation_date
    FROM %s.list l
    LEFT JOIN %s.user u ON u.id = l.user_id
    WHERE u.id = $1
  `, s.schema, s.schema)

	rows, err := s.session.Query(query, ownerID)
	defer rows.Close()
	if err != nil {
		return nil, domain.NewErrorWithCause(domain.UnknownErrorCode, err)
	}

	ownersLists := make([]*lists.List, 0)
	for rows.Next() {
		if err = rows.Err(); err != nil {
			return nil, domain.NewErrorWithCause(domain.UnknownErrorCode, err)
		}

		list := &lists.List{}

		err = rows.Scan(&list.ID, &list.Name, &list.CreationDate)
		if err != nil {
			return nil, domain.NewErrorWithCause(domain.UnknownErrorCode, err)
		}

		ownersLists = append(ownersLists, list)
	}

	return ownersLists, nil
}

// GetList retrieves a list with its id
func (s *PostgresStorage) GetList(id string) (*lists.List, *domain.Error) {
	query := fmt.Sprintf(`
    SELECT l.id, l.name, u.id
    FROM %s.list l
    LEFT JOIN %s.user u
    ON l.user_id = u.id
    WHERE l.id = $1
  `, s.schema, s.schema)

	row := s.session.QueryRow(query, id)

	l := &lists.List{
		Owner: &users.User{},
	}

	err := row.Scan(&l.ID, &l.Name, &l.Owner.ID)

	if err == sql.ErrNoRows {
		return nil, domain.NewError(lists.ListNotFoundErrorCode)
	}

	if err != nil {
		return nil, domain.NewErrorWithCause(domain.UnknownErrorCode, err)
	}

	return l, nil
}

func (s *PostgresStorage) GetProductInList(productID string, listID string) (*lists.ProductInList, *domain.Error) {
	query := fmt.Sprintf(`
    SELECT product_id, list_id
    FROM %s.product_in_list
    WHERE product_id = $1
    AND list_id = $2
  `, s.schema)

	row := s.session.QueryRow(query, productID, listID)

	pil := &lists.ProductInList{}

	err := row.Scan(&pil.ProductID, &pil.ListID)
	if err == sql.ErrNoRows {
		return nil, domain.NewError(lists.ProductInListNotFoundErrorCode)
	}
	if err != nil {
		return nil, domain.NewErrorWithCause(domain.UnknownErrorCode, err)
	}

	return pil, nil
}

func (s *PostgresStorage) AddProductToList(productID string, listID string) *domain.Error {
	query := fmt.Sprintf(`
    INSERT INTO %s.product_in_list(product_id, list_id)
    VALUES ($1, $2)
  `, s.schema)

	_, err := s.session.Exec(query, productID, listID)

	if err != nil {
		return domain.NewErrorWithCause(domain.UnknownErrorCode, err)
	}

	return nil
}

func (s *PostgresStorage) DeleteProductsFromList(listID string) *domain.Error {
	deleteProductInListQuery := fmt.Sprintf(`
    DELETE FROM %s.product_in_list WHERE list_id = $1
  `, s.schema)

	_, err := s.session.Exec(deleteProductInListQuery, listID)

	if err != nil {
		return domain.NewErrorWithCause(domain.UnknownErrorCode, err)
	}

	return nil
}

func (s *PostgresStorage) DeleteList(listID string) *domain.Error {
	deleteProductInListQuery := fmt.Sprintf(`
    DELETE FROM %s.list WHERE id = $1
  `, s.schema)

	_, err := s.session.Exec(deleteProductInListQuery, listID)

	if err != nil {
		return domain.NewErrorWithCause(domain.UnknownErrorCode, err)
	}

	return nil
}
