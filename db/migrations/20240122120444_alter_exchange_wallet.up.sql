ALTER TABLE exchange_wallet
ADD COLUMN last_used_in_block INT,
ADD COLUMN reference TEXT;
