package postgres

import (
	"fmt"
	"strings"

	"bitbucket.org/hornetdefiant/core-engine/factory"
	"github.com/spf13/viper"
	"golang.org/x/exp/slog"
)

func (p *Postgres) AddAssets(
	datas []factory.MonitorAssetsData,
	email string,
) error {
	if len(datas) == 0 {
		return nil
	}

	// Prepare the SQL statement for insertion
	query := "INSERT INTO monitor_assets (asset_type, asset_value,email) VALUES "
	values := []interface{}{}
	placeholders := []string{}

	for i, asset := range datas {
		placeholders = append(placeholders, fmt.Sprintf("($%d, $%d,$%d)", i*3+1, i*3+2, i*3+3))
		values = append(values, asset.AssetType, asset.AssetValue, email)
	}

	query += strings.Join(placeholders, ", ")
	stmt, err := p.dbConn.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute the SQL statement with provided data
	_, err = stmt.Exec(values...)
	return err
}

func (p *Postgres) GetAssets(pageNum int, email string) ([]factory.MonitorAssetsData, error) {
	ast := make([]factory.MonitorAssetsData, 0)
	args := make([]interface{}, 0)
	slog.Info("Running in Get Assets function in Postgres")
	pageSize := viper.GetInt("page_size")
	offset := (pageNum - 1) * pageSize
	args = append(args, email)
	args = append(args, pageSize)
	args = append(args, offset)
	query := "SELECT asset_id,asset_type,asset_value,add_date,risk_score from monitor_assets WHERE email = $1 LIMIT $2 OFFSET $3"
	slog.Info("query running")
	rows, err := p.dbConn.Query(query, args...)
	if err != nil {
		return ast, err
	}
	defer rows.Close()
	for rows.Next() {
		var t factory.MonitorAssetsData
		err = rows.Scan(
			&t.AssetId,
			&t.AssetType,
			&t.AssetValue,
			&t.Timestamp,
			&t.RiskScore,
		)
		if err != nil {
			return ast, err
		}
		ast = append(ast, t)
	}
	return ast, nil
}

func (p *Postgres) GetUserAssets() (map[string][]factory.AssetDetails, error) {
	query := `SELECT asset_id,email,asset_type,asset_value from monitor_assets`
	rows, err := p.dbConn.Query(query)
	if err != nil {
		return nil, err
	}
	ua := make(map[string][]factory.AssetDetails)
	var email string
	var ad factory.AssetDetails

	for rows.Next() {
		err := rows.Scan(&ad.AssetID, &email, &ad.AssetType, &ad.AssetValue)
		if err != nil {
			return nil, err
		}
		if asset, ok := ua[email]; ok {
			ua[email] = append(asset, ad)
		} else {
			uds := make([]factory.AssetDetails, 1)
			uds[0] = ad
			ua[email] = uds
		}
	}

	defer rows.Close()
	return ua, nil
}
