package contributors

type Repository interface {
	AddContributor(listUUID, userUUID string) error
	GetUserListUUIDs(userUUID string) ([]string, error)
}
