package factory

import "time"

type Tx struct {
	BlockNum  int       `json:"block_num"`
	Sender    string    `json:"sender"`
	Receiver  string    `json:"receiver"`
	Timestamp time.Time `json:"timestamp"`
	Amount    float64   `json:"amount"`
}
