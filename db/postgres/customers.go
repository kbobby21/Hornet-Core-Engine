package postgres

import (
	_ "github.com/lib/pq"
)

func (p *Postgres) CustomersAdd(email string) error {
	_, err := p.dbConn.Exec("INSERT INTO CUSTOMERS(email) VALUES($1)", email)
	return err
}
