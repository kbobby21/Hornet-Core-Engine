package factory

type DarkWebSite struct {
	WebsiteName  string          `json:"website_name"`
	OnionUrl     string          `json:"onion_url"`
	Wallets      []FlaggedWallet `json:"wallets"`
	Tags         []string        `json:"tags"`
	Body         []string        `json:"body"`
	DiscoveredAt string          `json:"created_at"`
}

type FlaggedWallet struct {
	Address    string `json:"address"`
	CryptoType string `json:"crypto_type"`
}

type DarkwebResponse struct {
	WebsiteName   string `json:"website_name"`
	OnionUrl      string `json:"onion_url"`
	WalletAddress string `json:"wallet_address"`
	Tag           string `json:"tag"`
	Body          string `json:"body"`
	DiscoveredAt  string `json:"discovered_at"`
}
