CREATE TABLE IF NOT EXISTS tenants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_id UUID NOT NULL REFERENCES customers(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'active',
    is_active BOOLEAN NOT NULL DEFAULT true,
    max_operators INTEGER NOT NULL DEFAULT 10,
    max_workcenters INTEGER NOT NULL DEFAULT 5,
    max_shop_floors INTEGER NOT NULL DEFAULT 3,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_tenants_customer_id ON tenants(customer_id);
CREATE INDEX idx_tenants_status ON tenants(status);
CREATE INDEX idx_tenants_is_active ON tenants(is_active);
