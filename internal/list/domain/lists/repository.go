package lists

type Repository interface {
	AddList(list *List) error
}
