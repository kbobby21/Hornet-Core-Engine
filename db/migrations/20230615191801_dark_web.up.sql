
CREATE TABLE IF NOT EXISTS dark_web_sites (
  id BIGSERIAL PRIMARY KEY,
  website_name TEXT NOT NULL,
  onion_url TEXT NOT NULL,
  wallet_address TEXT REFERENCES flagged_wallets(address),
  tag TEXT REFERENCES tags(value),
  body TEXT,
  discovered_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP 
);
