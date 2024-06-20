package postgres

import (
	"fmt"

	"bitbucket.org/hornetdefiant/core-engine/factory"
	"github.com/spf13/viper"
)

func (p *Postgres) AddNotification(notifications []factory.NotificationsData) error {
	query := "INSERT INTO notifications (asset_id, email, alert_type, alert_message, seen, nt_time) values"
	args := make([]interface{}, 0)
	for i, n := range notifications {
		query += fmt.Sprintf(
			"($%d, $%d, $%d, $%d, $%d, $%d)",
			6*i+1,
			6*i+2,
			6*i+3,
			6*i+4,
			6*i+5,
			6*i+6,
		)
		if i < len(notifications)-1 {
			query += ","
		}
		args = append(args,
			n.AssetId,
			n.Email,
			n.AlertType,
			n.AlertMessage,
			n.Seen,
			n.Timestamp,
		)
	}
	_, err := p.dbConn.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("error in inserting to notifications: %s", err.Error())
	}
	return nil
}

func (p *Postgres) GetNotifications(pageNum int, email string) ([]factory.NotificationsData, error) {
	notifications := make([]factory.NotificationsData, 0)
	args := make([]interface{}, 0)
	pageSize := viper.GetInt("page_size")
	offset := (pageNum - 1) * pageSize
	args = append(args, email)
	args = append(args, pageSize)
	args = append(args, offset)

	query := "SELECT email, alert_type, alert_message, seen, nt_time FROM notifications"
	conjunction := " WHERE email = $1 AND seen = false"
	order := " ORDER BY nt_time DESC LIMIT $2 OFFSET $3"

	query += conjunction + order

	rows, err := p.dbConn.Query(query, args...)
	if err != nil {
		return notifications, err
	}
	defer rows.Close()
	for rows.Next() {
		var n factory.NotificationsData
		err = rows.Scan(
			&n.Email,
			&n.AlertType,
			&n.AlertMessage,
			&n.Seen,
			&n.Timestamp,
		)
		if err != nil {
			return notifications, err
		}
		notifications = append(notifications, n)
	}
	return notifications, nil
}

func (p *Postgres) DeleteNotifications(email string) error {
	_, err := p.dbConn.Exec("UPDATE notifications SET seen = true WHERE email = $1", email)
	return err
}
