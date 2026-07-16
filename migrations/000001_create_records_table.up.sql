CREATE TABLE records (
    id SERIAL PRIMARY KEY,
    url TEXT NOT NULL,
    short_url VARCHAR(20) NOT NULL UNIQUE,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE INDEX idx_records_short_url ON records (short_url);