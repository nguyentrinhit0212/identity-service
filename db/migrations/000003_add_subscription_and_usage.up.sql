-- Add subscription fields
ALTER TABLE tenants
ADD COLUMN subscription_status VARCHAR(50),
ADD COLUMN subscription_plan VARCHAR(50),
ADD COLUMN subscription_expires_at TIMESTAMP WITH TIME ZONE,
ADD COLUMN usage_stats JSONB DEFAULT '{}'::jsonb;

-- Create index for subscription expiration
CREATE INDEX idx_tenant_subscription_expires ON tenants (subscription_expires_at)
WHERE subscription_expires_at IS NOT NULL;

-- Create index for subscription status
CREATE INDEX idx_tenant_subscription_status ON tenants (subscription_status)
WHERE subscription_status IS NOT NULL;

-- Update existing tenants with default values
UPDATE tenants
SET subscription_status = CASE
    WHEN type = 'personal' THEN 'free'
    WHEN type = 'team' THEN 'active'
    WHEN type = 'enterprise' THEN 'active'
    END,
    subscription_plan = type::text,
    subscription_expires_at = CASE
    WHEN type = 'personal' THEN NULL
    ELSE CURRENT_TIMESTAMP + INTERVAL '1 year'
    END,
    usage_stats = '{
        "apiCalls": 0,
        "storageUsed": 0,
        "lastUpdated": null
    }'::jsonb
WHERE subscription_status IS NULL; 