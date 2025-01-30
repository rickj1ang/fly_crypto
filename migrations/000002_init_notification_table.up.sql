-- Create notifications table
CREATE TABLE IF NOT EXISTS notifications (
    notification_id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(user_id),
    coin_symbol VARCHAR(20) NOT NULL,
    target_price DECIMAL(20, 8) NOT NULL,
    is_above BOOLEAN NOT NULL
);