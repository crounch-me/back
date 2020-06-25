package postgres

import (
	"database/sql"
	"fmt"

	"github.com/crounch-me/back/domain"
	"github.com/crounch-me/back/domain/products"
	"github.com/crounch-me/back/domain/users"
)

// CreateProduct inserts a new product
func (s *PostgresStorage) CreateProduct(id, name, ownerID string) *domain.Error {
	query := fmt.Sprintf(`
		INSERT INTO %s."product"(id, name, user_id)
		VALUES ($1, $2, $3)
	`, s.schema)

	_, err := s.session.Exec(query, id, name, ownerID)

	if err != nil {
		return domain.NewErrorWithCause(domain.UnknownErrorCode, err)
	}

	return nil
}

// GetProduct fetchs an existing product or return error
func (s *PostgresStorage) GetProduct(id string) (*products.Product, *domain.Error) {
	query := fmt.Sprintf(`
    SELECT p.id, p.name, u.id
    FROM %s.product p
    LEFT JOIN %s.user u
    ON p.user_id = u.id
    WHERE p.id = $1
  `, s.schema, s.schema)

	row := s.session.QueryRow(query, id)

	p := &products.Product{
		Owner: &users.User{},
	}

	err := row.Scan(&p.ID, &p.Name, &p.Owner.ID)

	if err == sql.ErrNoRows {
		return nil, domain.NewError(products.ProductNotFoundErrorCode)
	}
	if err != nil {
		return nil, domain.NewErrorWithCause(domain.UnknownErrorCode, err)
	}

	return p, nil
}
