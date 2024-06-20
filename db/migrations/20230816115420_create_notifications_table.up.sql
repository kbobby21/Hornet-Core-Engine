CREATE TABLE IF NOT EXISTS notifications (
    asset_id BIGSERIAL,
    email TEXT,
    alert_type TEXT,
    alert_message TEXT,
    seen    BOOLEAN DEFAULT FALSE,
    nt_time TIMESTAMPTZ
);