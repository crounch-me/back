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

// GetOwnerLists get all owner's lists
func (s *PostgresStorage) GetOwnerLists(ownerID string) ([]*model.List, error) {
	log.WithField("id", ownerID).Debug("Get lists of owner")

	query := fmt.Sprintf(`
    SELECT l.id, l.name
    FROM %s.list l
    LEFT JOIN %s.user u ON u.id = l.user_id
    WHERE u.id = $1
  `, s.schema, s.schema)

	rows, err := s.session.Query(query, ownerID)
	defer rows.Close()
	if err != nil {
		log.WithError(err).Error("Unable to get owner's list")
		return nil, err
	}

	lists := make([]*model.List, 0)
	for rows.Next() {
		if err = rows.Err(); err != nil {
			return nil, err
		}

		list := &model.List{}
		err = rows.Scan(&list.ID, &list.Name)
		if err != nil {
			return nil, err
		}
		lists = append(lists, list)
	}
	return lists, nil
}
