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
	contributors    []*Contributor
	products        []*Product
}

func NewList(uuid string, name string, creationDate time.Time, archivationDate *time.Time) (*List, error) {
	if uuid == "" {
		return nil, errors.New("empty list uuid")
	}

	if name == "" {
		return nil, errors.New("empty list name")
	}

	if creationDate.IsZero() {
		return nil, errors.New("empty list creation date")
	}

	if archivationDate != nil && archivationDate.IsZero() {
		return nil, errors.New("empty but not nil list archivation date")
	}

	return &List{
		uuid:            uuid,
		name:            name,
		creationDate:    creationDate,
		archivationDate: archivationDate,
		products:        []*Product{},
		contributors:    []*Contributor{},
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

func (l *List) IsArchived() bool {
	return l.archivationDate != nil && !l.archivationDate.IsZero()
}

func (l *List) Archive() {
	now := time.Now()

	l.archivationDate = &now
}
