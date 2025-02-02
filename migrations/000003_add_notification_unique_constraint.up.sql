-- Add unique constraint for coin_symbol and is_above combination
ALTER TABLE notifications
ADD CONSTRAINT unique_user_coin_above UNIQUE (user_id, coin_symbol, is_above);