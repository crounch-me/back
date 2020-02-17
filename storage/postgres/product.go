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
