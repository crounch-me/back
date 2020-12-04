package lists

import "errors"

type Contributor struct {
	uuid string
}

func NewContributor(uuid string) (*Contributor, error) {
	if uuid == "" {
		return nil, errors.New("contributor uuid is empty")
	}

	return &Contributor{uuid: uuid}, nil
}

func (c Contributor) UUID() string {
	return c.uuid
}

func (l *List) AddContributor(c *Contributor) {
	l.contributors = append(l.contributors, c)
}

func (l List) Contributors() []*Contributor {
	return l.contributors
}

func (c Contributor) ContributeIn(l *List) bool {
	for _, listContributor := range l.Contributors() {
		if listContributor.uuid == c.uuid {
			return true
		}
	}

	return false
}
