package factory

type BeneficiarySummary struct {
	Receiver      string  `json:"receiver"`
	TotalReceived float64 `json:"total_received"`
}
