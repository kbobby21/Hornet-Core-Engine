package beneficiary

import "bitbucket.org/hornetdefiant/core-engine/factory"

type Repository interface {
	GetBeneficiary(sender string) ([]factory.BeneficiarySummary, error)
}
