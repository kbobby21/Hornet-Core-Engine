CREATE TABLE exchange_contact (
  exchange_id BIGINT REFERENCES exchange_metadata(id),
  contact_email TEXT,
  contact_number TEXT
);