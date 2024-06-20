package users

import "bitbucket.org/hornetdefiant/core-engine/factory"

type Repository interface {
	AddUser(data factory.User) error
	LoginUser(data factory.User) (string, bool, error)
	VerifyEmail(email string) error
	GetAllUsersEmail() ([]string, error)
	GetUser(email string) (*factory.User, error)
}
