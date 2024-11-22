CREATE TABLE ip_whitelist (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    ip_address INET UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);
