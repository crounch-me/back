package postgres

import (
	"fmt"

	"github.com/Sehsyha/crounch-back/model"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

// CreateProduct inserts a new product
func (s *PostgresStorage) CreateProduct(product *model.Product) error {
	log.WithField("name", product.Name).Debug("Creating product")

	product.ID = uuid.NewV4().String()
	query := fmt.Sprintf(`
		INSERT INTO %s."product"(id, name, user_id)
		VALUES ($1, $2, $3)
	`, s.schema)

	_, err := s.session.Exec(query, product.ID, product.Name, product.Owner.ID)

	if err != nil {
		log.WithError(err).Error("Unable to create product")
		return err
	}

	return nil
}

// GetProduct fetchs an existing product or return error
func (s *PostgresStorage) GetProduct(id string) (*model.Product, error) {
	log.WithField("userID", id).Debug("Getting product")

	query := fmt.Sprintf(`
    SELECT p.id, p.name, u.id
    FROM %s.product p
    LEFT JOIN %s.user u
    ON p.user_id = u.id
    WHERE p.id = $1
  `, s.schema, s.schema)

	row := s.session.QueryRow(query, id)

	p := &model.Product{
		Owner: &model.User{},
	}

	err := row.Scan(&p.ID, &p.Name, &p.Owner.ID)
	err = handleNotFound(err)

	if err != nil {
		return nil, err
	}

	return p, nil
}
