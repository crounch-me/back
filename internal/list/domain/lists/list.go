package lists

import (
	"errors"
	"time"
)

type List struct {
	uuid            string
	name            string
	creationDate    time.Time
	archivationDate *time.Time
	contributors    []string
	products        []*Product
}

var (
	ErrProductAlreadyInList = errors.New("product already in lists")
)

func NewList(uuid string, name string, creatorUUID string, creationDate time.Time) (*List, error) {
	if uuid == "" {
		return nil, errors.New("empty list uuid")
	}

	if name == "" {
		return nil, errors.New("empty list name")
	}

	if creatorUUID == "" {
		return nil, errors.New("empty list creator uuid")
	}

	if creationDate.IsZero() {
		return nil, errors.New("empty list creation date")
	}

	return &List{
		uuid:         uuid,
		name:         name,
		creationDate: creationDate,
		products:     []*Product{},
	}, nil
}

func (l List) UUID() string {
	return l.uuid
}

func (l List) Name() string {
	return l.name
}

func (l List) CreationDate() time.Time {
	return l.creationDate
}

func (l List) HasProduct(uuid string) bool {
	for _, p := range l.products {
		if p.uuid == uuid {
			return true
		}
	}

	return false
}
