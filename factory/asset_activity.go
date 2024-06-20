package factory

type AssetActivity struct {
	AssetID    int      `json:"asset_id"`
	Activities []string `json:"activities"`
}

type AssetDetails struct {
	AssetID    int    `json:"asset_id"`
	AssetType  string `json:"asset_type"`
	AssetValue string `json:"asset_value"`
}
