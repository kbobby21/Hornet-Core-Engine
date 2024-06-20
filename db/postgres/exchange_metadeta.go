package postgres

import (
	"fmt"
	"strings"

	"bitbucket.org/hornetdefiant/core-engine/factory"
	"bitbucket.org/hornetdefiant/core-engine/utils"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

func (p *Postgres) AddExchangeMetaData(data *factory.ExchangeMetaDataInsert) error {

	// Begin a transaction
	tx, err := p.dbConn.Begin()
	if err != nil {
		return err
	}

	// Insert data into the exchange_metadata table
	var exchangeID int64
	err = tx.QueryRow("INSERT INTO exchange_metadata (name, country) VALUES ($1, $2) RETURNING id", data.Name, data.Country).Scan(&exchangeID)
	if err != nil {
		tx.Rollback()
		return err
	}

	contactEmail, contactNumber := utils.BigLength(data.ContactEmail, data.ContactNumber)
	for i := 0; i < len(contactEmail); i++ {
		// Insert data into the exchange_contact table
		_, err = tx.Exec("INSERT INTO exchange_contact (exchange_id, contact_email, contact_number) VALUES ($1, $2, $3)", exchangeID, strings.TrimSpace(contactEmail[i]), strings.TrimSpace(contactNumber[i]))
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// Insert data into the exchange_wallet table
	for _, v := range utils.CommaSeparateArr(data.WalletAddress) {
		_, err = tx.Exec("INSERT INTO exchange_wallet (exchange_id, wallet_address, last_used_in_block, references) VALUES ($1, $2, $3, $4)", exchangeID, strings.TrimSpace(v))
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		// return err
		return err
	}
	return nil
}

func (p *Postgres) GetExchangeMetaData(pageNum int) ([]factory.Exchange, error) {
	exchanges := make([]factory.Exchange, 0)
	pageSize := viper.GetInt("page_size")
	offset := (pageNum - 1) * pageSize
	query := `
        SELECT em.name, em.country, ec.contact_email, ec.contact_number, ew.wallet_address, ew.last_used_in_block, ew.reference
        FROM exchange_metadata em
        JOIN exchange_contact ec ON ec.exchange_id = em.id
        JOIN exchange_wallet ew ON ew.exchange_id = em.id
		LIMIT $1 OFFSET $2
    `

	// Execute the query and fetch the results
	rows, err := p.dbConn.Query(query, pageSize, offset) //db.Query(query)
	if err != nil {
		return exchanges, err
	}
	defer rows.Close()

	// Loop over the rows and populate the Exchange structs
	for rows.Next() {
		exchange := factory.Exchange{}
		if err := rows.Scan(&exchange.Name, &exchange.Country, &exchange.ContactEmail, &exchange.ContactNumber, &exchange.WalletAddress, &exchange.LastUsedInBlock, &exchange.Reference); err != nil {
			return exchanges, err
		}
		exchanges = append(exchanges, exchange)
	}
	fmt.Printf("exchanges: %+v\n", exchanges)
	return exchanges, nil
}
