-- Add tenant_id column to users table
ALTER TABLE users 
ADD COLUMN tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE;

-- Create index for tenant_id
CREATE INDEX idx_users_tenant_id ON users(tenant_id);
