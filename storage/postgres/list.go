package postgres

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/crounch-me/back/internal/categories"
	"github.com/crounch-me/back/internal/common/errors"
	"github.com/crounch-me/back/internal/listing"
	"github.com/crounch-me/back/internal/products"
	"github.com/crounch-me/back/util"
)

// CreateList inserts a new list
func (s *PostgresStorage) CreateList(id, name string, creationDate time.Time) *errors.Error {
	query := fmt.Sprintf(`
		INSERT INTO %s."list"(id, name, creation_date)
		VALUES ($1, $2, $3)
	`, s.schema)

	_, err := s.session.Exec(query, id, name, creationDate)

	if err != nil {
		return errors.NewError(errors.UnknownErrorCode).WithCause(err)
	}

	return nil
}

// ArchiveList sets the archivation date into the list table
func (s *PostgresStorage) ArchiveList(listID string, archivationDate time.Time) *errors.Error {
	query := fmt.Sprintf(`
    UPDATE %s.list
    SET archivation_date = $1
    WHERE id = $2
  `, s.schema)

	_, err := s.session.Exec(query, archivationDate, listID)
	if err != nil {
		return errors.NewError(errors.UnknownErrorCode).WithCause(err)
	}

	return nil
}

// GetUsersLists get all user's lists
func (s *PostgresStorage) GetUsersLists(userID string) ([]*listing.List, *errors.Error) {
	query := fmt.Sprintf(`
    SELECT l.id, l.name, l.creation_date, l.archivation_date
    FROM %s.list l
    LEFT JOIN %s.contributor c ON c.list_id = l.id
    WHERE c.user_id = $1
  `, s.schema, s.schema)

	rows, err := s.session.Query(query, userID)
	defer rows.Close()
	if err != nil {
		return nil, errors.NewError(errors.UnknownErrorCode).WithCause(err)
	}

	ownersLists := make([]*listing.List, 0)
	for rows.Next() {
		if err = rows.Err(); err != nil {
			return nil, errors.NewError(errors.UnknownErrorCode).WithCause(err)
		}

		list := &listing.List{}

		err = rows.Scan(&list.ID, &list.Name, &list.CreationDate, &list.ArchivationDate)
		if err != nil {
			return nil, errors.NewError(errors.UnknownErrorCode).WithCause(err)
		}

		ownersLists = append(ownersLists, list)
	}

	return ownersLists, nil
}

// GetList retrieves a list with its id
func (s *PostgresStorage) GetList(id string) (*listing.List, *errors.Error) {
	query := fmt.Sprintf(`
    SELECT l.id, l.name, l.creation_date, l.archivation_date
    FROM %s.list l
    WHERE l.id = $1
  `, s.schema)

	callError := &errors.CallError{
		PackageName: "postgres",
		MethodName:  "GetList",
	}

	row := s.session.QueryRow(query, id)

	l := &listing.List{}
	err := row.Scan(&l.ID, &l.Name, &l.CreationDate, &l.ArchivationDate)

	if err == sql.ErrNoRows {
		return nil, errors.NewError(listing.ListNotFoundErrorCode).WithCallError(callError)
	}

	if err != nil {
		return nil, errors.NewError(errors.UnknownErrorCode).WithCallError(callError).WithCause(err)
	}

	return l, nil
}

// UpdateProductInList updates the bought value in product in list
func (s *PostgresStorage) UpdateProductInList(updateProductInList *listing.UpdateProductInList, productID, listID string) (*listing.ProductInListLink, *errors.Error) {
	query := fmt.Sprintf(`
    UPDATE %s.product_in_list
    SET bought = $1
    WHERE product_id = $2
    AND list_id = $3
    RETURNING product_id, list_id, bought
  `, s.schema)

	row := s.session.QueryRow(query, updateProductInList.Bought, productID, listID)

	pil := &listing.ProductInListLink{}

	err := row.Scan(&pil.ProductID, &pil.ListID, &pil.Bought)
	if err == sql.ErrNoRows {
		logger := util.GetLogger()
		logger.WithError(err).
			WithField("package", "postgres").
			Debug("UpdateProductInList")
		return nil, errors.NewError(listing.ProductInListNotFoundErrorCode)
	}

	return pil, nil
}

func (s *PostgresStorage) GetProductInList(productID string, listID string) (*listing.ProductInListLink, *errors.Error) {
	query := fmt.Sprintf(`
    SELECT product_id, list_id
    FROM %s.product_in_list
    WHERE product_id = $1
    AND list_id = $2
  `, s.schema)

	row := s.session.QueryRow(query, productID, listID)

	pil := &listing.ProductInListLink{}

	err := row.Scan(&pil.ProductID, &pil.ListID)
	if err == sql.ErrNoRows {
		return nil, errors.NewError(listing.ProductInListNotFoundErrorCode)
	}
	if err != nil {
		return nil, errors.NewError(errors.UnknownErrorCode).WithCause(err)
	}

	return pil, nil
}

func (s *PostgresStorage) GetProductsOfList(listID string) ([]*listing.ProductInList, *errors.Error) {
	query := fmt.Sprintf(`
    SELECT p.id, p.name, pil.bought, c.id, c.name
    FROM %s.product p
    LEFT JOIN %s.product_in_list pil ON pil.product_id = p.id
    LEFT JOIN %s.list l ON pil.list_id = l.id
    LEFT JOIN %s.category c ON c.id = p.category_id
    WHERE l.id = $1
  `, s.schema, s.schema, s.schema, s.schema)

	rows, err := s.session.Query(query, listID)
	defer rows.Close()
	if err != nil {
		return nil, errors.NewError(errors.UnknownErrorCode).WithCause(err)
	}

	productsOfList := make([]*listing.ProductInList, 0)
	for rows.Next() {
		if err = rows.Err(); err != nil {
			return nil, errors.NewError(errors.UnknownErrorCode).WithCause(err)
		}

		productOfList := &listing.ProductInList{
			Product: &products.Product{},
		}
		var nullableCategoryID, nullableCategoryName sql.NullString

		err = rows.Scan(&productOfList.Product.ID, &productOfList.Product.Name, &productOfList.Bought, &nullableCategoryID, &nullableCategoryName)
		if err != nil {
			return nil, errors.NewError(errors.UnknownErrorCode).WithCause(err)
		}

		if nullableCategoryID.Valid && nullableCategoryName.Valid {
			productOfList.Product.Category = &categories.Category{
				ID:   nullableCategoryID.String,
				Name: nullableCategoryName.String,
			}
		}

		productsOfList = append(productsOfList, productOfList)
	}

	return productsOfList, nil
}

func (s *PostgresStorage) AddProductToList(productID string, listID string) *errors.Error {
	query := fmt.Sprintf(`
    INSERT INTO %s.product_in_list(product_id, list_id)
    VALUES ($1, $2)
  `, s.schema)

	_, err := s.session.Exec(query, productID, listID)
	if err != nil {
		return errors.NewError(errors.UnknownErrorCode).WithCause(err)
	}

	return nil
}

func (s *PostgresStorage) DeleteProductFromList(productID string, listID string) *errors.Error {
	query := fmt.Sprintf(`
    DELETE FROM %s.product_in_list
    WHERE product_id = $1
    AND list_id = $2
  `, s.schema)

	_, err := s.session.Exec(query, productID, listID)
	if err != nil {
		return errors.NewError(errors.UnknownErrorCode).WithCause(err)
	}

	return nil
}

func (s *PostgresStorage) DeleteProductsFromList(listID string) *errors.Error {
	deleteProductInListQuery := fmt.Sprintf(`
    DELETE FROM %s.product_in_list WHERE list_id = $1
  `, s.schema)

	_, err := s.session.Exec(deleteProductInListQuery, listID)
	if err != nil {
		return errors.NewError(errors.UnknownErrorCode).WithCause(err)
	}

	return nil
}

func (s *PostgresStorage) DeleteList(listID string) *errors.Error {
	deleteProductInListQuery := fmt.Sprintf(`
    DELETE FROM %s.list WHERE id = $1
  `, s.schema)

	_, err := s.session.Exec(deleteProductInListQuery, listID)
	if err != nil {
		return errors.NewError(errors.UnknownErrorCode).WithCause(err)
	}

	return nil
}
