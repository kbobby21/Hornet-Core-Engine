CREATE TABLE exchange_wallet (
  exchange_id BIGINT REFERENCES exchange_metadata(id),
  wallet_address TEXT
);