package lists

type Repository interface {
	AddList(list *List) error
	ReadByIDs(uuids []string) ([]*List, error)
	ReadByID(uuid string) (*List, error)
	AddContributor(listUUID, contributorUUID string) error
	GetContributorListUUIDs(contributorUUID string) ([]string, error)
}
