-- Create users table
CREATE TABLE IF NOT EXISTS users (
    user_id BIGSERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    notifications_id BIGINT[] DEFAULT ARRAY[]::BIGINT[]
);