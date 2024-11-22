CREATE TABLE failed_logins (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL,
    ip_address VARCHAR(45) NOT NULL,
    failed_at TIMESTAMP DEFAULT NOW(),
    attempt_count INTEGER DEFAULT 1,
    reason VARCHAR(255)
);
