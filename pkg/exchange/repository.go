package exchange

import "bitbucket.org/hornetdefiant/core-engine/factory"

type Repository interface {
	AddExchangeMetaData(data *factory.ExchangeMetaDataInsert) error
	GetExchangeMetaData(pageNum int) ([]factory.Exchange, error)
}
