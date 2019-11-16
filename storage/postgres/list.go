package postgres

import (
	"fmt"

	"github.com/Sehsyha/crounch-back/model"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

// CreateList inserts a new list
func (s *PostgresStorage) CreateList(list *model.List) error {
	log.WithField("name", list.Name).Debug("Creating list")

	list.ID = uuid.NewV4().String()
	query := fmt.Sprintf(`
		INSERT INTO %s."list"(id, name, user_id)
		VALUES ($1, $2, $3)
	`, s.schema)

	_, err := s.session.Exec(query, list.ID, list.Name, list.Owner.ID)

	if err != nil {
		log.WithError(err).Error("Unable to create list")
		return err
	}

	return nil
}
