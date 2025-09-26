-- Drop users table and related objects
-- This reverses the changes made in 000001_create_users_table.up.sql

-- Drop trigger
DROP TRIGGER IF EXISTS update_users_updated_at ON users;

-- Drop trigger function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop indexes (they will be dropped automatically with the table, but explicit for clarity)
DROP INDEX IF EXISTS idx_users_updated_at;
DROP INDEX IF EXISTS idx_users_created_at;
DROP INDEX IF EXISTS idx_users_is_deleted;
DROP INDEX IF EXISTS idx_users_deleted_at;
DROP INDEX IF EXISTS idx_email_active;
DROP INDEX IF EXISTS idx_username_active;

-- Drop users table
DROP TABLE IF EXISTS users;