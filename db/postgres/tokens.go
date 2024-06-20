package postgres

import (
	"database/sql"
	"errors"
	"strings"

	"bitbucket.org/hornetdefiant/core-engine/factory"
	"golang.org/x/exp/slog"
)

func (p *Postgres) AddToken(
	data factory.Tokens,
	email string,
) error {
	var admin = "@hornet.technology"
	if strings.Contains(email, admin) {
		data.IsAdmin = true
	} else {
		data.IsAdmin = false
	}
	query := "INSERT INTO tokens (email, token, is_admin, privileges, valid_till) VALUES ($1, $2, $3, $4, $5) ON CONFLICT (email) DO UPDATE SET token = EXCLUDED.token, privileges = EXCLUDED.privileges, valid_till = EXCLUDED.valid_till"
	_, err := p.dbConn.Exec(query, email, data.Token, data.IsAdmin, data.Privileges, data.ValidTill)
	return err
}

func (p *Postgres) CheckAdminAndValidToken(token string) (bool, string, error) {
	slog.Info("Running token validation in postgres")
	var isAdmin bool
	var email string

	query := "SELECT email, valid_till > Now() AND is_admin FROM tokens WHERE token = $1"
	err := p.dbConn.QueryRow(query, token).Scan(&email, &isAdmin)
	if err != nil {
		if err == sql.ErrNoRows {
			slog.Info("No rows found for the token in postgres")
			return false, "", errors.New("please generate a token for this email")
		}
		slog.Info("Error in token validation in postgres: ", err.Error())
		return false, "", err
	}

	return isAdmin, email, nil
}

func (p *Postgres) CheckAdminAndValidTokenByEmail(email string) (bool, error) {
	slog.Info("Running token validation by email in postgres")

	if strings.Contains(email, "@hornet.technology") {
		return true, nil
	}
	var isAdmin bool
	query := "SELECT valid_till > Now() AND is_admin FROM tokens WHERE email = $1"
	err := p.dbConn.QueryRow(query, email).Scan(&isAdmin)
	if err != nil {
		if err == sql.ErrNoRows {
			slog.Info("No rows found for the email in postgres")
			return false, errors.New("please generate a token for this email")
		}
		slog.Info("Error in token validation by email in postgres: ", err.Error())
		return false, err
	}

	return isAdmin, nil
}
