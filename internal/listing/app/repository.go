package app

import "github.com/crounch-me/back/internal/listing/domain/lists"

type Repository interface {
	DeleteList(uuid string) error
	ReadByContributor(c *lists.Contributor) ([]*lists.List, error)
	ReadByIDs(uuids []string) ([]*lists.List, error)
	ReadByID(uuid string) (*lists.List, error)
	SaveList(l *lists.List) error
	UpdateList(l *lists.List) error
}
