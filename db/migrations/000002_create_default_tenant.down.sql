-- Drop trigger and function
DROP TRIGGER IF EXISTS create_personal_tenant_trigger ON users;
DROP FUNCTION IF EXISTS create_personal_tenant();

-- Drop indexes
DROP INDEX IF EXISTS idx_personal_tenant_owner;
DROP INDEX IF EXISTS idx_tenant_domain;

-- Drop columns
ALTER TABLE tenants DROP COLUMN IF EXISTS type;
ALTER TABLE tenants DROP COLUMN IF EXISTS owner_id;
ALTER TABLE tenants DROP COLUMN IF EXISTS max_users;
ALTER TABLE tenants DROP COLUMN IF EXISTS features;
ALTER TABLE tenants DROP COLUMN IF EXISTS domain_verified;

-- Drop enum type
DROP TYPE IF EXISTS tenant_type; 