package postgres

import (
	"database/sql"
	"errors"

	"bitbucket.org/hornetdefiant/core-engine/factory"
)

func (p *Postgres) GetTransactionInfo(traddress, coaddress string) ([]factory.TransactionInfo, error) {
	var direction string
	var amount int
	var blockNum int
	var txtime sql.NullTime
	var counterParty string

	query := `
        SELECT 
            CASE
                WHEN sender = $1 AND receiver = $2 THEN 'sent to'
                WHEN sender = $2 AND receiver = $1 THEN 'received from'
                ELSE 'unknown'
            END AS direction,
            CASE
                WHEN sender = $1 AND receiver = $2 THEN amount
                WHEN sender = $2 AND receiver = $1 THEN amount
                ELSE 0
            END AS amount,
            CASE
                WHEN sender = $1 AND receiver = $2 THEN block_num
                WHEN sender = $2 AND receiver = $1 THEN block_num
                ELSE 0
            END AS block_num,
            CASE
                WHEN sender = $1 AND receiver = $2 THEN txtime
                WHEN sender = $2 AND receiver = $1 THEN txtime
                ELSE NULL
            END AS txtime,
            CASE
                WHEN sender = $1 AND receiver = $2 THEN $2
                WHEN sender = $2 AND receiver = $1 THEN $1
                ELSE ''
            END AS counter_party
        FROM txs
        WHERE (sender = $1 AND receiver = $2) OR (sender = $2 AND receiver = $1);
    `

	err := p.dbConn.QueryRow(query, traddress, coaddress).Scan(&direction, &amount, &blockNum, &txtime, &counterParty)
	if err == sql.ErrNoRows {
		// Handle the case when no rows were found
		return nil, errors.New("No transaction found for the given addresses")
	} else if err != nil {
		// Handle other errors
		return nil, err
	}

	transaction := factory.TransactionInfo{
		Direction:    direction,
		Amount:       amount,
		BlockNum:     blockNum,
		CounterParty: counterParty,
	}

	if txtime.Valid {
		transaction.Txtime = txtime.Time
	}

	return []factory.TransactionInfo{transaction}, nil
}
