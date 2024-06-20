package tags

import "bitbucket.org/hornetdefiant/core-engine/factory"

type Repository interface {
	GetTags() ([]factory.Tag, error)
}
