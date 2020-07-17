package postgres

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/crounch-me/back/domain"
	"github.com/crounch-me/back/domain/categories"
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
		return domain.NewError(domain.UnknownErrorCode).WithCause(err)
	}

	return nil
}

// GetProduct fetchs an existing product or return error
func (s *PostgresStorage) GetProduct(id string) (*products.Product, *domain.Error) {
	query := fmt.Sprintf(`
    SELECT p.id, p.name, u.id, c.id, c.name
    FROM %s.product p
    LEFT JOIN %s.user u ON p.user_id = u.id
    LEFT JOIN %s.category c ON p.category_id = c.id
    WHERE p.id = $1
  `, s.schema, s.schema, s.schema)

	row := s.session.QueryRow(query, id)

	p := &products.Product{
		Owner: &users.User{},
	}

	var nullableCategoryID, nullableCategoryName sql.NullString

	err := row.Scan(&p.ID, &p.Name, &p.Owner.ID, &nullableCategoryID, &nullableCategoryName)

	if err == sql.ErrNoRows {
		return nil, domain.NewError(products.ProductNotFoundErrorCode)
	}
	if err != nil {
		return nil, domain.NewError(domain.UnknownErrorCode)
	}

	if nullableCategoryName.Valid && nullableCategoryID.Valid {
		p.Category = &categories.Category{
			ID:   nullableCategoryID.String,
			Name: nullableCategoryName.String,
		}
	}

	return p, nil
}

func (s *PostgresStorage) SearchDefaults(name string, userID string) ([]*products.Product, *domain.Error) {
	lowerCasedName := strings.ToLower(name)

	query := fmt.Sprintf(`
    SELECT p.id as product_id, p.name as product_name, c.id as category_id, c.name as category_name
    FROM %s.product p
    LEFT JOIN %s.user u ON u.id = p.user_id
    LEFT JOIN %s.category c ON c.id = p.category_id
    WHERE LOWER(p.name) SIMILAR TO '(' || $1 || '%%|%% ' || $1 || '%%)'
    AND u.id = $2
  `, s.schema, s.schema, s.schema)

	rows, err := s.session.Query(query, lowerCasedName, userID)
	if err != nil {
		return nil, domain.NewError(domain.UnknownErrorCode).WithCause(err)
	}
	defer rows.Close()

	productList := make([]*products.Product, 0)
	for rows.Next() {
		if err = rows.Err(); err != nil {
			return nil, domain.NewError(domain.UnknownErrorCode).WithCause(err)
		}

		product := &products.Product{}
		var nullableCategoryID, nullableCategoryName sql.NullString

		err = rows.Scan(&product.ID, &product.Name, &nullableCategoryID, &nullableCategoryName)
		if err != nil {
			return nil, domain.NewError(domain.UnknownErrorCode).WithCause(err)
		}

		if nullableCategoryName.Valid && nullableCategoryID.Valid {
			product.Category = &categories.Category{
				ID:   nullableCategoryID.String,
				Name: nullableCategoryName.String,
			}
		}

		productList = append(productList, product)
	}

	return productList, nil
}
