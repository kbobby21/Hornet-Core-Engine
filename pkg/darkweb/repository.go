package darkweb

import "bitbucket.org/hornetdefiant/core-engine/factory"

type Repository interface {
	AddWebSiteMeta(data factory.DarkWebSite) error
	GetWebsiteMeta(pageNo int) ([]factory.DarkwebResponse, error)
}
