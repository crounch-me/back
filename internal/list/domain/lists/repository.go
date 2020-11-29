package lists

type Repository interface {
	AddList(list *List) error
	ReadByIDs(uuids []string) ([]*List, error)
}
