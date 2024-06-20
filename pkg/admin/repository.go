package admin

type Repository interface {
	CheckAdminAndValidToken(token string) (bool, string, error)
	CheckAdminAndValidTokenByEmail(email string) (bool, error)
}
