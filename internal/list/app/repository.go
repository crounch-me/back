package app

import "github.com/crounch-me/back/internal/list/domain/lists"

type Repository interface {
	ReadByContributor(c *lists.Contributor) ([]*lists.List, error)
	ReadByIDs(uuids []string) ([]*lists.List, error)
	ReadByID(uuid string) (*lists.List, error)
	SaveList(l *lists.List) error
	UpdateList(l *lists.List) error
}
