CREATE TABLE payments (
    id SERIAL PRIMARY KEY,
    order_id INT REFERENCES orders(id) ON DELETE CASCADE,
    provider TEXT,
    amount NUMERIC(10,2),
    status TEXT CHECK (status IN ('pending','success','failed')) DEFAULT 'pending',
    transaction_id TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);
