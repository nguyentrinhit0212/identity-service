-- Drop indexes
DROP INDEX IF EXISTS idx_sessions_user_id;
DROP INDEX IF EXISTS idx_sessions_tenant_id;
DROP INDEX IF EXISTS idx_sessions_expires_at;

-- Drop table
DROP TABLE IF EXISTS sessions; 