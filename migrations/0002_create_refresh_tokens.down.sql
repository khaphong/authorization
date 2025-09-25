-- Migration: 0002_create_refresh_tokens.down.sql
-- Description: Drop refresh tokens table and related objects

-- Drop indexes
DROP INDEX IF EXISTS idx_refresh_tokens_token_hash;
DROP INDEX IF EXISTS idx_refresh_tokens_user_id;
DROP INDEX IF EXISTS idx_refresh_tokens_expires_at;
DROP INDEX IF EXISTS idx_refresh_tokens_revoked;
DROP INDEX IF EXISTS idx_refresh_tokens_created_at;
DROP INDEX IF EXISTS idx_refresh_tokens_active;

-- Drop table
DROP TABLE IF EXISTS refresh_tokens;
