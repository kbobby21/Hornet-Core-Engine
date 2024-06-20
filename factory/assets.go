package factory

import "time"

type MonitorAssetsData struct {
	AssetId    int       `json:"asset_id"`
	AssetType  string    `json:"asset_type"`
	AssetValue string    `json:"asset_value"`
	Timestamp  time.Time `json:"timestamp"`
	RiskScore  float32   `json:"risk_score"`
}
