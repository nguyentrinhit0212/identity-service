-- Remove usage_stats column from tenants table
ALTER TABLE tenants DROP COLUMN IF EXISTS usage_stats; 