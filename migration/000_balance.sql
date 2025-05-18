-- Create Balance Table
CREATE TABLE IF NOT EXISTS balances (
                                        id SERIAL PRIMARY KEY,
                                        user_id INT UNIQUE NOT NULL,
                                        balance FLOAT DEFAULT 0,
                                        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);