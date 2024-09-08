CREATE TABLE IF NOT EXISTS users (
    uuid UUID NOT NULL,
    national_id TEXT PRIMARY KEY,
    email TEXT,
    last_name TEXT,
    ip TEXT,
    state TEXT
);