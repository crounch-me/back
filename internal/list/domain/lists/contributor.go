package lists

import "errors"

var (
	ErrContributorAlreadyInList = errors.New("contributor already in list")
)

func (l List) HasContributor(uuid string) bool {
	for _, c := range l.contributors {
		if c == uuid {
			return true
		}
	}

	return false
}

func (l *List) AddContributor(uuid string) error {
	if l.HasContributor(uuid) {
		return ErrContributorAlreadyInList
	}

	l.contributors = append(l.contributors, uuid)

	return nil
}
