package contributors

type Repository interface {
	AddContributor(listUUID, userUUID string) error
}
