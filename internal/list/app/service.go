package app

import (
	"errors"
	"time"

	"github.com/crounch-me/back/internal/list/domain/lists"
	"github.com/crounch-me/back/util"
)

type ListService struct {
	listsRepository Repository
}

func NewListService(listRepository Repository) (*ListService, error) {
	if listRepository == nil {
		return nil, errors.New("listRepository is nil")
	}

	return &ListService{
		listsRepository: listRepository,
	}, nil
}

func (s *ListService) CreateList(creatorUUID, name string) (string, error) {
	listUUID, err := util.GenerateID()
	if err != nil {
		return "", err
	}

	creationDate := time.Now()

	l, err := lists.NewList(listUUID, name, creationDate, nil)
	if err != nil {
		return "", err
	}

	c, err := lists.NewContributor(creatorUUID)
	if err != nil {
		return "", err
	}

	l.AddContributor(c)

	err = s.listsRepository.SaveList(l)
	if err != nil {
		return "", err
	}

	return listUUID, nil
}

func (s *ListService) GetContributorLists(contributorUUID string) ([]*lists.List, error) {
	c, err := lists.NewContributor(contributorUUID)
	if err != nil {
		return nil, err
	}

	return s.listsRepository.ReadByContributor(c)
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

func (s *ListService) ArchiveList(listUUID string) (*lists.List, error) {
	l, err := s.listsRepository.ReadByID(listUUID)
	if err != nil {
		return nil, err
	}

	l.Archive()

	err = s.listsRepository.UpdateList(l)
	if err != nil {
		return nil, err
	}

	return l, nil
}
