-- Add usage_stats column to tenants table
ALTER TABLE tenants ADD COLUMN usage_stats JSONB DEFAULT '{}'::jsonb; 