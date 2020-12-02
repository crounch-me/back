package app

import (
	"errors"
	"time"

	"github.com/crounch-me/back/internal/list/domain/lists"
	"github.com/crounch-me/back/util"
)

type ListService struct {
	listsRepository lists.Repository
}

func NewListService(listRepository lists.Repository) (*ListService, error) {
	if listRepository == nil {
		return nil, errors.New("listRepository is nil")
	}

	return &ListService{
		listsRepository: listRepository,
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

	err = s.listsRepository.AddContributor(listUUID, userUUID)
	if err != nil {
		return "", err
	}

	return listUUID, nil
}

func (s *ListService) GetUserLists(userUUID string) ([]*lists.List, error) {
	listUUIDs, err := s.listsRepository.GetContributorListUUIDs(userUUID)
	if err != nil {
		return nil, err
	}

	return s.listsRepository.ReadByIDs(listUUIDs)
}

func (s *ListService) ReadList(userUUID, listUUID string) (*lists.List, error) {
	l, err := s.listsRepository.ReadByID(listUUID)
	if err != nil {
		return nil, err
	}

	c, err := lists.NewContributor(userUUID)
	if err != nil {
		return nil, err
	}

	if !c.ContributeIn(l) {
		return nil, errors.New("user not contributor")
	}

	return l, nil
}
