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

func (c *Category) UUID() string {
	return c.uuid
}

func (c *Category) Name() string {
	return c.name
}

func (p *Product) HasCategory() bool {
	return p.category != nil
}

func (p *Product) CategoryName() string {
	return p.category.Name()
}

func (p *Product) CategoryUUID() string {
	return p.category.UUID()
}
