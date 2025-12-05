-- Add tenant_id column to customers table
ALTER TABLE customers 
ADD COLUMN tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE;

-- Create index for tenant_id
CREATE INDEX idx_customers_tenant_id ON customers(tenant_id);
