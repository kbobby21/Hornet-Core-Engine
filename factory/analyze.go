package factory

import "time"

type TransactionInfo struct {
	Direction    string    `json:"direction"`
	Amount       int       `json:"amount"`
	BlockNum     int       `json:"block_num"`
	Txtime       time.Time `json:"txtime"`
	CounterParty string    `json:"counter_party"`
}
