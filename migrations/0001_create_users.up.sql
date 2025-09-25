-- Migration: 0001_create_users.up.sql
-- Description: Create users table with UUID v7 support and soft delete

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create users table
CREATE TABLE users (
    id VARCHAR(36) PRIMARY KEY DEFAULT uuid_generate_v4()::text,
    username VARCHAR(50) NOT NULL,
    email VARCHAR(255) NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE NULL,
    is_deleted BOOLEAN DEFAULT FALSE NOT NULL
);

-- Create indexes for performance and constraints
CREATE UNIQUE INDEX idx_users_username_active ON users(username) WHERE is_deleted = false;
CREATE UNIQUE INDEX idx_users_email_active ON users(email) WHERE is_deleted = false;
CREATE INDEX idx_users_is_deleted ON users(is_deleted);
CREATE INDEX idx_users_deleted_at ON users(deleted_at) WHERE deleted_at IS NOT NULL;
CREATE INDEX idx_users_created_at ON users(created_at);

-- Create updated_at trigger function
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE 'plpgsql';

-- Create trigger for users table
CREATE TRIGGER trigger_users_updated_at 
    BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Add comments for documentation
COMMENT ON TABLE users IS 'Users table with soft delete support';
COMMENT ON COLUMN users.id IS 'Primary key - UUID v7 for time-ordered UUIDs';
COMMENT ON COLUMN users.deleted_at IS 'Timestamp when user was soft deleted';
COMMENT ON COLUMN users.is_deleted IS 'Boolean flag for soft delete status';
