package authorizations

type Repository interface {
	AddAuthorization(userUUID, token string) error
	GetUserUUIDByToken(token string) (string, error)
	RemoveAuthorization(userUUID, token string) error
}
