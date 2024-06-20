CREATE TABLE flagged_wallets (
	address TEXT PRIMARY KEY,
        crypto_type TEXT NOT NULL,
        discovered_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
