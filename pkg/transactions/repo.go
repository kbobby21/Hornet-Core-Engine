package transactions

import "bitbucket.org/hornetdefiant/core-engine/factory"

type Repo interface {
	AddTransactions([]factory.Tx) error
	GetTransactions(
		pageNum int,
		sender,
		reciever string,
	) ([]factory.Tx, error)
}
