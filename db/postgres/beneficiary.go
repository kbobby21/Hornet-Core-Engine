package postgres

import (
	"bitbucket.org/hornetdefiant/core-engine/factory"
)

func (p *Postgres) GetBeneficiary(sender string) ([]factory.BeneficiarySummary, error) {
	var beneficiarySummaries []factory.BeneficiarySummary

	query := `
        SELECT receiver, SUM(amount) as total_received
        FROM txs
        WHERE sender = $1
        GROUP BY receiver
    `

	rows, err := p.dbConn.Query(query, sender)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var beneficiarySummary factory.BeneficiarySummary
		err := rows.Scan(&beneficiarySummary.Receiver, &beneficiarySummary.TotalReceived)
		if err != nil {
			return nil, err
		}
		beneficiarySummaries = append(beneficiarySummaries, beneficiarySummary)
	}

	return beneficiarySummaries, nil
}
