DROP INDEX IF EXISTS idx_jwt_tokens_tenant_id;
ALTER TABLE jwt_tokens DROP COLUMN tenant_id; 