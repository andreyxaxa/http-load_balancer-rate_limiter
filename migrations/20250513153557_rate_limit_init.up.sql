CREATE TABLE IF NOT EXISTS rate_limits (
    client_ip  VARCHAR PRIMARY KEY NOT NULL, 
    capacity   INTEGER NOT NULL CHECK (capacity > 0),
    fill_rate  DOUBLE PRECISION NOT NULL CHECK (fill_rate >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);