ALTER TABLE jwt_tokens
ADD COLUMN tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE;
CREATE INDEX idx_jwt_tokens_tenant_id ON jwt_tokens(tenant_id); 