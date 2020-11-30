package app

import (
	"errors"
	"time"

	"github.com/crounch-me/back/internal/list/domain/contributors"
	"github.com/crounch-me/back/internal/list/domain/lists"
	"github.com/crounch-me/back/util"
)

type ListService struct {
	listsRepository        lists.Repository
	contributorsRepository contributors.Repository
}

func NewListService(listRepository lists.Repository, contributorsRepository contributors.Repository) (*ListService, error) {
	if listRepository == nil {
		return nil, errors.New("listRepository is nil")
	}

	if contributorsRepository == nil {
		return nil, errors.New("contributorsRepository is nil")
	}

	return &ListService{
		listsRepository:        listRepository,
		contributorsRepository: contributorsRepository,
	}, nil
}

func (s *ListService) CreateList(userUUID, name string) (string, error) {
	listUUID, err := util.GenerateID()
	if err != nil {
		return "", err
	}

	creationDate := time.Now()

	list, err := lists.NewList(listUUID, name, creationDate, nil)
	if err != nil {
		return "", err
	}

	err = s.listsRepository.AddList(list)
	if err != nil {
		return "", err
	}

	err = s.contributorsRepository.AddContributor(listUUID, userUUID)
	if err != nil {
		return "", err
	}

	return listUUID, nil
}

func (s *ListService) GetUserLists(userUUID string) ([]*lists.List, error) {
	listUUIDs, err := s.contributorsRepository.GetUserListUUIDs(userUUID)
	if err != nil {
		return nil, err
	}

	return s.listsRepository.ReadByIDs(listUUIDs)
}
