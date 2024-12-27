-- Drop triggers
DROP TRIGGER IF EXISTS update_user_tenant_access_updated_at ON user_tenant_access;
DROP TRIGGER IF EXISTS update_tenants_updated_at ON tenants;
DROP TRIGGER IF EXISTS update_user_credentials_updated_at ON user_credentials;
DROP TRIGGER IF EXISTS update_users_updated_at ON users;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop tables (in correct order due to foreign key constraints)
DROP TABLE IF EXISTS failed_logins;
DROP TABLE IF EXISTS ip_whitelists;
DROP TABLE IF EXISTS oauth_providers;
DROP TABLE IF EXISTS jwt_tokens;
DROP TABLE IF EXISTS user_tenant_access;
DROP TABLE IF EXISTS user_credentials;
DROP TABLE IF EXISTS tenants;
DROP TABLE IF EXISTS users;

-- Drop extension if no other tables need it
DROP EXTENSION IF EXISTS "uuid-ossp"; 