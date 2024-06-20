package assets

import "bitbucket.org/hornetdefiant/core-engine/factory"

type Repository interface {
	AddAssets(datas []factory.MonitorAssetsData, email string) error
	GetAssets(page_no int, email string) ([]factory.MonitorAssetsData, error)
	GetUserAssets() (map[string][]factory.AssetDetails, error)
}
