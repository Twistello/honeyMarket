CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    total NUMERIC(10,2),
    status TEXT CHECK (status IN ('new','paid','shipped','completed','cancelled')) DEFAULT 'new',
    created_at TIMESTAMP DEFAULT NOW()
);
