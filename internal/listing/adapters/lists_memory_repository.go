package adapters

import (
	"errors"
	"sync"

	"github.com/crounch-me/back/internal/listing/app"
	"github.com/crounch-me/back/internal/listing/domain/lists"
)

type ListsMemoryRepository struct {
	lists map[string]*lists.List
	lock  *sync.RWMutex
}

func NewListsMemoryRepository() app.Repository {
	return &ListsMemoryRepository{
		lists: make(map[string]*lists.List, 0),
		lock:  &sync.RWMutex{},
	}
}

func (r *ListsMemoryRepository) DeleteList(uuid string) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	if _, ok := r.lists[uuid]; !ok {
		return errors.New("list not found")
	}

	delete(r.lists, uuid)

	return nil
}

func (r *ListsMemoryRepository) ReadByContributor(c *lists.Contributor) ([]*lists.List, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	lists := make([]*lists.List, 0)
	for _, l := range r.lists {
		if c.ContributeIn(l) {
			lists = append(lists, l)
		}
	}

	return lists, nil
}

func (r *ListsMemoryRepository) ReadByIDs(uuids []string) ([]*lists.List, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	lists := make([]*lists.List, 0)
	for _, uuid := range uuids {
		if l, ok := r.lists[uuid]; ok {
			lists = append(lists, l)
		}
	}

	return lists, nil
}

func (r *ListsMemoryRepository) ReadByID(uuid string) (*lists.List, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	l, ok := r.lists[uuid]
	if !ok {
		return nil, lists.ErrListNotFound
	}

	return l, nil
}

func (r *ListsMemoryRepository) SaveList(l *lists.List) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	if _, ok := r.lists[l.UUID()]; ok {
		return errors.New("list already exists")
	}

	r.lists[l.UUID()] = l

	return nil
}

func (r *ListsMemoryRepository) UpdateList(l *lists.List) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	if _, ok := r.lists[l.UUID()]; !ok {
		return errors.New("list not found")
	}

	r.lists[l.UUID()] = l

	return nil
}
