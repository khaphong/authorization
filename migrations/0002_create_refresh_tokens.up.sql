-- Migration: 0002_create_refresh_tokens.up.sql
-- Description: Create refresh tokens table

CREATE TABLE refresh_tokens (
    id VARCHAR(36) PRIMARY KEY DEFAULT uuid_generate_v4()::text,
    user_id VARCHAR(36) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash VARCHAR(64) NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    revoked BOOLEAN DEFAULT FALSE NOT NULL
);

-- Create indexes for performance
CREATE UNIQUE INDEX idx_refresh_tokens_token_hash ON refresh_tokens(token_hash);
CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens(user_id);
CREATE INDEX idx_refresh_tokens_expires_at ON refresh_tokens(expires_at);
CREATE INDEX idx_refresh_tokens_revoked ON refresh_tokens(revoked);
CREATE INDEX idx_refresh_tokens_created_at ON refresh_tokens(created_at);

-- Composite index for common queries
CREATE INDEX idx_refresh_tokens_active ON refresh_tokens(token_hash, revoked, expires_at) 
    WHERE revoked = false;

-- Add comments for documentation
COMMENT ON TABLE refresh_tokens IS 'Refresh tokens for JWT authentication';
COMMENT ON COLUMN refresh_tokens.token_hash IS 'SHA-256 hash of the refresh token';
COMMENT ON COLUMN refresh_tokens.revoked IS 'Boolean flag indicating if token is revoked';
