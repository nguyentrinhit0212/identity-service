-- Remove subscription_plan column from tenants table
ALTER TABLE tenants DROP COLUMN IF EXISTS subscription_plan; 