package factory

type ExchangeMetaDataInsert struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Country       string `json:"country"`
	ContactEmail  string `json:"contact_email"`
	ContactNumber string `json:"contact_number"`
	WalletAddress string `json:"wallet_address"`
	// ExchangeContact ExchangeContact `json:"exchange_contact"`
	// ExchangeWallet  ExchangeWallet  `json:"exchange_wallet"`
}

type ExchangeContact struct {
	ID            int    `json:"id"`
	ExchangeId    string `json:"-"`
	ContactEmail  string `json:"contact_email"`
	ContactNumber string `json:"contact_number"`
}

type ExchangeWallet struct {
	ID              int    `json:"id"`
	ExchangeId      string `json:"-"`
	WalletAddress   string `json:"wallet_address"`
	LastUsedInBlock int    `json:"last_used_in_block"`
	Reference       string `json:"reference"`
}

type Exchange struct {
	Name            string `json:"name"`
	Country         string `json:"country"`
	ContactEmail    string `json:"contact_email"`
	ContactNumber   string `json:"contact_number"`
	WalletAddress   string `json:"wallet_address"`
	LastUsedInBlock int    `json:"last_used_in_block"`
	Reference       string `json:"reference"`
}
