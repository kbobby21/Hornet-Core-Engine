package postgres

import (
	"bitbucket.org/hornetdefiant/core-engine/factory"
	"github.com/lib/pq"
	"github.com/spf13/viper"
)

func (n *Postgres) GetAssetActivities(pageNum int, email string) ([]factory.AssetActivity, error) {

	activities := make([]factory.AssetActivity, 0)
	pageSize := viper.GetInt("page_size")
	offset := (pageNum - 1) * pageSize

	query := `
        SELECT
            asset_id,
            array_agg(alert_message) AS activities
        FROM
            notifications
        WHERE
            email = $1
        GROUP BY
            asset_id
        ORDER BY
            asset_id
        LIMIT $2 OFFSET $3
    `

	rows, err := n.dbConn.Query(query, email, pageSize, offset)
	if err != nil {
		return activities, err
	}
	defer rows.Close()

	for rows.Next() {
		var activity factory.AssetActivity
		var activityArray []string
		err = rows.Scan(
			&activity.AssetID,
			pq.Array(&activityArray),
		)
		if err != nil {
			return activities, err
		}
		activity.Activities = activityArray
		activities = append(activities, activity)
	}
	return activities, nil
}
