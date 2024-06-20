CREATE TABLE IF NOT EXISTS users (
    email varchar(100) PRIMARY KEY,
    password TEXT,
    verified BOOLEAN DEFAULT FALSE
);