package postgres

import (
	"database/sql"
	"errors"

	"bitbucket.org/hornetdefiant/core-engine/factory"
	_ "github.com/lib/pq"
)

func (p *Postgres) AddUser(data factory.User) error {
	var organization string

	if data.Organization == "" {
		organization = "Individual"
	} else {
		organization = data.Organization
	}

	_, err := p.dbConn.Exec("INSERT INTO users (email, password, organization) VALUES ($1, $2, $3)", data.Email, data.HashedPassword, organization)
	return err
}

// LoginUser returns the password and status of veriied (bool) for the given user
func (p *Postgres) LoginUser(data factory.User) (string, bool, error) {
	var pass string
	var verified bool
	err := p.dbConn.QueryRow("select password,verified from users  where email = $1", data.Email).Scan(&pass, &verified)

	return pass, verified, err
}

func (p *Postgres) VerifyEmail(email string) error {
	_, err := p.dbConn.Exec(
		"UPDATE users SET verified=TRUE where email=$1",
		email,
	)
	return err
}

func (p *Postgres) GetAllUsersEmail() ([]string, error) {
	rows, err := p.dbConn.Query("SELECT email FROM users WHERE email <> ''")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var emails []string
	for rows.Next() {
		var email string
		if err := rows.Scan(&email); err != nil {
			return nil, err
		}
		emails = append(emails, email)
	}

	return emails, nil
}

func (p *Postgres) GetUser(email string) (*factory.User, error) {
	query := `SELECT u.email, u.organization, t.token,t.valid_till
    FROM users u
    LEFT JOIN tokens t ON u.email = t.email
    WHERE u.email = $1;`

	rows, err := p.dbConn.Query(query, email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var user factory.User
	var token sql.NullString
	var validTill sql.NullTime

	if rows.Next() {
		err := rows.Scan(&user.Email, &user.Organization, &token,&validTill)
		if err != nil {
			return nil, err
		}

		user.Token = token.String
		user.ValidTill = validTill.Time
		return &user, nil
	}

	return nil, errors.New("user not found")
}
