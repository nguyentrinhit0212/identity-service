-- Drop indexes
DROP INDEX IF EXISTS idx_tenant_subscription_expires;
DROP INDEX IF EXISTS idx_tenant_subscription_status;

-- Drop subscription and usage columns
ALTER TABLE tenants
DROP COLUMN IF EXISTS subscription_status,
DROP COLUMN IF EXISTS subscription_plan,
DROP COLUMN IF EXISTS subscription_expires_at,
DROP COLUMN IF EXISTS usage_stats; 