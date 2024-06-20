package token

import "bitbucket.org/hornetdefiant/core-engine/factory"

type Repository interface {
	AddToken(data factory.Tokens, email string) error
}
