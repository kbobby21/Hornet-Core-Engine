package analyze

import "bitbucket.org/hornetdefiant/core-engine/factory"

type Repository interface {
	GetTransactionInfo(traddress, coaddress string) ([]factory.TransactionInfo, error)
}
