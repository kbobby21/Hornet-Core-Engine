CREATE TABLE IF NOT EXISTS monitor_assets (
    asset_id BIGSERIAL,
    add_date TIMESTAMPTZ DEFAULT NOW(), 
    email TEXT,
    asset_type TEXT NOT NULL,
    asset_value TEXT NOT NULL,
    risk_score DECIMAL DEFAULT 0.0,
    PRIMARY KEY(email,asset_value)
);
