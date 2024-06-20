package postgres

import (
	"fmt"

	"github.com/spf13/viper"
	"bitbucket.org/hornetdefiant/core-engine/factory"
)

func (p *Postgres) AddTransactions(txs []factory.Tx) error {
	query := "INSERT INTO txs(block_num,sender,receiver,txtime,amount) values "
	args := make([]interface{}, 0)
	for i, tx := range txs {
		query += fmt.Sprintf(
			"($%d,$%d,$%d,$%d,$%d)",
			5*i+1,
			5*i+2,
			5*i+3,
			5*i+4,
			5*i+5,
		)
		if i < len(txs)-1 {
			query += ","
		}
		args = append(
			args,
			tx.BlockNum,
			tx.Sender,
			tx.Receiver,
			tx.Timestamp,
			tx.Amount,
		)
	}
	_, err := p.dbConn.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("error in inserting to txs: %s", err.Error())
	}
	return nil
}

func (p *Postgres) GetTransactions(
	pageNum int,
	sender string,
	receiver string,
) ([]factory.Tx, error) {
	txs := make([]factory.Tx, 0)
	args := make([]interface{}, 0)

	pageSize := viper.GetInt("page_size")
	offset := (pageNum - 1) * pageSize

	args = append(args, pageSize)
	args = append(args, offset)
	query := "SELECT block_num,sender,receiver,txtime,amount from txs "
	order := "ORDER BY txtime DESC LIMIT $1 OFFSET $2"

	conjunction := "WHERE "

	i := 3 // because first two args are already added to args variable
	if sender != "all" {
		query += conjunction + fmt.Sprintf("sender=$%d ", i)
		i++
		args = append(args, sender)
		conjunction = "AND "
	}
	if receiver != "all" {
		query += conjunction + fmt.Sprintf("receiver=$%d ", i)
		args = append(args, receiver)
	}

	query += order
	rows, err := p.dbConn.Query(query, args...)
	if err != nil {
		return txs, err
	}
	defer rows.Close()

	for rows.Next() {
		var t factory.Tx
		err = rows.Scan(
			&t.BlockNum,
			&t.Sender,
			&t.Receiver,
			&t.Timestamp,
			&t.Amount,
		)
		if err != nil {
			return txs, err
		}
		txs = append(txs, t)
	}
	return txs, nil
}
