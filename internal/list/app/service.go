package app

import (
	"time"

	"github.com/crounch-me/back/internal/list/domain/lists"
	"github.com/crounch-me/back/util"
)

type ListService struct {
	repository lists.Repository
}

func NewListService(listRepository lists.Repository) *ListService {
	return &ListService{
		repository: listRepository,
	}
}

func (l *ListService) CreateList(userUUID, name string) (string, error) {
	listUUID, err := util.GenerateID()
	if err != nil {
		return "", err
	}

	creationDate := time.Now()

	list, err := lists.NewList(listUUID, name, userUUID, creationDate)
	if err != nil {
		return "", err
	}

	err = l.repository.AddList(list)
	if err != nil {
		return "", err
	}

	return listUUID, nil
}
