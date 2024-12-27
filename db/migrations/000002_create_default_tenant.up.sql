-- Create tenant type enum
CREATE TYPE tenant_type AS ENUM ('personal', 'team', 'enterprise');

-- Add tenant type and related columns
ALTER TABLE tenants ADD COLUMN type tenant_type NOT NULL DEFAULT 'team';
ALTER TABLE tenants ADD COLUMN owner_id UUID REFERENCES users(id);
ALTER TABLE tenants ADD COLUMN max_users INTEGER;
ALTER TABLE tenants ADD COLUMN features JSONB DEFAULT '{}'::jsonb;
ALTER TABLE tenants ADD COLUMN subscription_status VARCHAR(50);
ALTER TABLE tenants ADD COLUMN subscription_expires_at TIMESTAMP WITH TIME ZONE;
ALTER TABLE tenants ADD COLUMN domain_verified BOOLEAN DEFAULT false;

-- Add constraints
CREATE UNIQUE INDEX idx_personal_tenant_owner ON tenants (owner_id) WHERE type = 'personal';
CREATE UNIQUE INDEX idx_tenant_domain ON tenants (domain) WHERE domain IS NOT NULL AND domain_verified = true;

-- Create function to automatically create personal tenant for new users
CREATE OR REPLACE FUNCTION create_personal_tenant()
RETURNS TRIGGER AS $$
BEGIN
    -- Create personal tenant
    INSERT INTO tenants (
        slug,
        name,
        type,
        owner_id,
        max_users,
        features,
        settings
    ) VALUES (
        'personal-' || NEW.id,  -- Unique slug based on user ID
        NEW.name || '''s Workspace',
        'personal',
        NEW.id,
        1,  -- Personal tenants limited to 1 user
        '{"personalFeatures": true}'::jsonb,
        '{"isPersonal": true}'::jsonb
    );

    -- Get the tenant ID we just created
    WITH new_tenant AS (
        SELECT id FROM tenants 
        WHERE owner_id = NEW.id AND type = 'personal'
        LIMIT 1
    )
    -- Grant admin access to the user for their personal tenant
    INSERT INTO user_tenant_access (
        user_id,
        tenant_id,
        roles,
        permissions
    )
    SELECT 
        NEW.id,
        new_tenant.id,
        ARRAY['admin']::text[],
        ARRAY['*']::text[]
    FROM new_tenant;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger to automatically create personal tenant for new users
CREATE TRIGGER create_personal_tenant_trigger
AFTER INSERT ON users
FOR EACH ROW
EXECUTE FUNCTION create_personal_tenant();

-- Insert default tenant template for each type
INSERT INTO tenants (
    slug,
    name,
    type,
    max_users,
    features,
    settings
) VALUES 
    (
        'default-personal',
        'Personal Workspace',
        'personal',
        1,
        '{"personalFeatures": true}'::jsonb,
        '{"isDefaultTemplate": true, "type": "personal"}'::jsonb
    ),
    (
        'default-team',
        'Team Workspace',
        'team',
        10,
        '{"teamFeatures": true, "collaborationTools": true}'::jsonb,
        '{"isDefaultTemplate": true, "type": "team"}'::jsonb
    ),
    (
        'default-enterprise',
        'Enterprise Workspace',
        'enterprise',
        NULL, -- Unlimited users
        '{"enterpriseFeatures": true, "sso": true, "audit": true, "advancedSecurity": true}'::jsonb,
        '{"isDefaultTemplate": true, "type": "enterprise"}'::jsonb
    )
ON CONFLICT (slug) DO NOTHING; 