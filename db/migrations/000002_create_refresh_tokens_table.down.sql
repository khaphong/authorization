-- Drop refresh_tokens table and related objects
-- This reverses the changes made in 000002_create_refresh_tokens_table.up.sql

-- Drop foreign key constraint first
ALTER TABLE refresh_tokens DROP CONSTRAINT IF EXISTS fk_refresh_tokens_user_id;

-- Drop indexes (they will be dropped automatically with the table, but explicit for clarity)
DROP INDEX IF EXISTS idx_refresh_tokens_created_at;
DROP INDEX IF EXISTS idx_refresh_tokens_revoked;
DROP INDEX IF EXISTS idx_refresh_tokens_expires_at;
DROP INDEX IF EXISTS idx_refresh_tokens_user_id;
DROP INDEX IF EXISTS idx_refresh_tokens_token_hash;

-- Drop refresh_tokens table
DROP TABLE IF EXISTS refresh_tokens;