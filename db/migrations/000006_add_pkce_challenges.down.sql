DROP INDEX IF EXISTS idx_pkce_challenges_used;
DROP INDEX IF EXISTS idx_pkce_challenges_expires_at;
DROP INDEX IF EXISTS idx_pkce_challenges_tenant_id;
DROP INDEX IF EXISTS idx_pkce_challenges_user_id;
DROP TABLE IF EXISTS pkce_challenges; 