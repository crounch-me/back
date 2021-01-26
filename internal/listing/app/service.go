package app

import (
	"errors"
	"time"

	"github.com/crounch-me/back/internal/common/utils"
	"github.com/crounch-me/back/internal/listing/domain/lists"
)

type ListService struct {
	listsRepository   Repository
	generationLibrary utils.GenerationLibrary
}

func NewListService(listRepository Repository, generationLibrary utils.GenerationLibrary) (*ListService, error) {
	if listRepository == nil {
		return nil, errors.New("listRepository is nil")
	}

	if generationLibrary == nil {
		return nil, errors.New("generationService is nil")
	}

	return &ListService{
		listsRepository:   listRepository,
		generationLibrary: generationLibrary,
	}, nil
}

func (s *ListService) CreateList(creatorUUID, name string) (string, error) {
	listUUID, err := s.generationLibrary.UUID()
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

func (s *ListService) ArchiveList(contributorUUID, listUUID string) error {
	l, err := s.listsRepository.ReadByID(listUUID)
	if err != nil {
		return err
	}

	c, err := lists.NewContributor(contributorUUID)
	if err != nil {
		return err
	}

	if !c.ContributeIn(l) {
		return errors.New("user not contributor")
	}

	l.Archive()

	err = s.listsRepository.UpdateList(l)
	if err != nil {
		return err
	}

	return nil
}

func (s *ListService) DeleteList(contributorUUID, listUUID string) error {
	l, err := s.listsRepository.ReadByID(listUUID)
	if err != nil {
		return err
	}

	c, err := lists.NewContributor(contributorUUID)
	if err != nil {
		return err
	}

	if !c.ContributeIn(l) {
		return errors.New("user not contributor")
	}

	err = s.listsRepository.DeleteList(l.UUID())
	if err != nil {
		return err
	}

	return nil
}

func (s *ListService) AddProductToList(contributorUUID, productUUID, listUUID string) error {
	l, err := s.listsRepository.ReadByID(listUUID)
	if err != nil {
		return err
	}

	c, err := lists.NewContributor(contributorUUID)
	if err != nil {
		return err
	}

	if !c.ContributeIn(l) {
		return errors.New("user not contributor")
	}

	p, err := lists.NewProduct(productUUID)
	if err != nil {
		return err
	}

	if l.HasProduct(p) {
		return errors.New("product already in list")
	}

	l.AddProduct(p)

	return s.listsRepository.UpdateList(l)
}

func (s *ListService) BuyProductInList(contributorUUID, productUUID, listUUID string) error {
	l, err := s.listsRepository.ReadByID(listUUID)
	if err != nil {
		return err
	}

	c, err := lists.NewContributor(contributorUUID)
	if err != nil {
		return err
	}

	if !c.ContributeIn(l) {
		return errors.New("user not contributor")
	}

	p, err := lists.NewProduct(productUUID)
	if err != nil {
		return err
	}

	return l.Buy(p)
}

func (s *ListService) DeleteProductInList(contributorUUID, productUUID, listUUID string) error {
	l, err := s.listsRepository.ReadByID(listUUID)
	if err != nil {
		return err
	}

	c, err := lists.NewContributor(contributorUUID)
	if err != nil {
		return err
	}

	if !c.ContributeIn(l) {
		return errors.New("user not contributor")
	}

	p, err := lists.NewProduct(productUUID)
	if err != nil {
		return err
	}

	return l.RemoveProduct(p)
}
