-- init-database.sql
-- This file ensures the database and required extensions are created

-- Create database if it doesn't exist (this will be handled by POSTGRES_DB env var)
-- But we can ensure extensions are available

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create a simple health check function
CREATE OR REPLACE FUNCTION health_check() 
RETURNS TEXT AS $$
BEGIN
    RETURN 'OK';
END;
$$ LANGUAGE plpgsql;
