CREATE TABLE IF NOT EXISTS txs (
	id BIGSERIAL PRIMARY KEY,
	block_num BIGINT, 
	sender TEXT,
	receiver TEXT,
	txtime TIMESTAMPTZ,
	amount real);

