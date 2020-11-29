package app

import (
	"time"

	"github.com/crounch-me/back/internal/list/domain/contributors"
	"github.com/crounch-me/back/internal/list/domain/lists"
	"github.com/crounch-me/back/util"
)

type ListService struct {
	listsRepository        lists.Repository
	contributorsRepository contributors.Repository
}

func NewListService(listRepository lists.Repository) *ListService {
	return &ListService{
		listsRepository: listRepository,
	}
}

func (l *ListService) CreateList(userUUID, name string) (string, error) {
	listUUID, err := util.GenerateID()
	if err != nil {
		return "", err
	}

	creationDate := time.Now()

	list, err := lists.NewList(listUUID, name, creationDate)
	if err != nil {
		return "", err
	}

	err = l.listsRepository.AddList(list)
	if err != nil {
		return "", err
	}

	err = l.contributorsRepository.AddContributor(listUUID, userUUID)
	if err != nil {
		return "", err
	}

	return listUUID, nil
}

func (l *ListService) GetUserLists(userUUID string) ([]*lists.List, error) {
	listUUIDs, err := l.contributorsRepository.GetUserListUUIDs(userUUID)
	if err != nil {
		return nil, err
	}

	return l.listsRepository.ReadByIDs(listUUIDs)
}
