CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    role TEXT CHECK (role IN ('customer','admin')) DEFAULT 'customer',
    created_at TIMESTAMP DEFAULT NOW()
);