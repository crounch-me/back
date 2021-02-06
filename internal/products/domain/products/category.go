package products

import "errors"

type Category struct {
	uuid string
	name string
}

func NewCategory(uuid, name string) (*Category, error) {
	if uuid == "" {
		return nil, errors.New("category uuid is empty")
	}

	if name == "" {
		return nil, errors.New("category name is empty")
	}

	return &Category{
		uuid: uuid,
		name: name,
	}, nil
}
