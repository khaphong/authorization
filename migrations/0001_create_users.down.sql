-- Migration: 0001_create_users.down.sql
-- Description: Drop users table and related objects

-- Drop trigger
DROP TRIGGER IF EXISTS trigger_users_updated_at ON users;

-- Drop trigger function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop indexes
DROP INDEX IF EXISTS idx_users_username_active;
DROP INDEX IF EXISTS idx_users_email_active;
DROP INDEX IF EXISTS idx_users_is_deleted;
DROP INDEX IF EXISTS idx_users_deleted_at;
DROP INDEX IF EXISTS idx_users_created_at;

-- Drop table
DROP TABLE IF EXISTS users;
