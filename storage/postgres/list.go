package postgres

import (
	"database/sql"
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

// GetList retrieves a list with its id
func (s *PostgresStorage) GetList(id string) (*model.List, error) {
	log.WithField("id", id).Debug("Get list")

	query := fmt.Sprintf(`
    SELECT l.id, l.name, u.id
    FROM %s.list l
    LEFT JOIN %s.user u
    ON l.user_id = u.id
    WHERE l.id = $1
  `, s.schema, s.schema)

	row := s.session.QueryRow(query, id)

	l := &model.List{
		Owner: &model.User{},
	}

	err := row.Scan(&l.ID, &l.Name, &l.Owner.ID)

	if err == sql.ErrNoRows {
		return nil, model.NewDatabaseError(model.ErrNotFound, nil)
	}

	if err != nil {
		return nil, err
	}

	return l, nil
}

func (s *PostgresStorage) GetOFFProducts(listID string) ([]*model.OFFProduct, error) {
	log.WithField("listID", listID).Debug("Get OFF products")

	query := fmt.Sprintf(`
    SELECT code
    FROM %s.list_off_product
    WHERE list_id = $1
  `, s.schema)

	rows, err := s.session.Query(query, listID)
	defer rows.Close()
	if err != nil {
		log.WithError(err).Error("Unable to get off products")
		return nil, err
	}

	products := make([]*model.OFFProduct, 0)
	for rows.Next() {
		if err = rows.Err(); err != nil {
			return nil, err
		}

		p := &model.OFFProduct{}
		err = rows.Scan(&p.Code)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}

// AddOFFProductToList adds a product from open food facts to an user list
func (s *PostgresStorage) AddOFFProductToList(listID string, offProduct *model.OFFProduct) error {
	log.WithField("listID", listID).WithField("offProduct.code", offProduct.Code).Debug("Add OFF product to list")

	query := fmt.Sprintf(`
		INSERT INTO %s."list_off_product"(list_id, code)
		VALUES ($1, $2)
	`, s.schema)

	_, err := s.session.Exec(query, listID, offProduct.Code)

	if err != nil {
		log.WithError(err).Error("Unable to add off product to list")
		return err
	}

	return nil
}
