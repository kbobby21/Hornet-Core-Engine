CREATE TABLE IF NOT EXISTS tokens (
    email TEXT,
    token TEXT,
    is_admin BOOLEAN,
    privileges TEXT,
    creation_time TIMESTAMPTZ DEFAULT NOW(),
    valid_till TIMESTAMP
);