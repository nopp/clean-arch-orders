
CREATE TABLE IF NOT EXISTS orders (
	id TEXT PRIMARY KEY,
	customer_name TEXT NOT NULL,
	total_amount NUMERIC(12,2) NOT NULL CHECK (total_amount >= 0),
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
