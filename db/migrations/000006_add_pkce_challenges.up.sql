CREATE TABLE IF NOT EXISTS pkce_challenges (
    id UUID PRIMARY KEY,
    code_challenge VARCHAR(255) NOT NULL,
    code_verifier VARCHAR(255) NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    used BOOLEAN NOT NULL DEFAULT FALSE
);

-- Add indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_pkce_challenges_user_id ON pkce_challenges(user_id);
CREATE INDEX IF NOT EXISTS idx_pkce_challenges_tenant_id ON pkce_challenges(tenant_id);
CREATE INDEX IF NOT EXISTS idx_pkce_challenges_expires_at ON pkce_challenges(expires_at);
CREATE INDEX IF NOT EXISTS idx_pkce_challenges_used ON pkce_challenges(used); 