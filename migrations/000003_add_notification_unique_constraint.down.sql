-- Remove unique constraint for coin_symbol and is_above combination
ALTER TABLE notifications
DROP CONSTRAINT IF EXISTS unique_user_coin_above;