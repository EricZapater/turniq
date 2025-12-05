-- Create operators table
CREATE TABLE IF NOT EXISTS operators (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    shop_floor_id UUID NOT NULL REFERENCES shopfloors(id) ON DELETE CASCADE,
    customer_id UUID NOT NULL REFERENCES customers(id) ON DELETE CASCADE,
    code VARCHAR(50) NOT NULL,
    name VARCHAR(255) NOT NULL,
    surname VARCHAR(255),
    vat_number VARCHAR(50),
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes
CREATE INDEX idx_operators_tenant_id ON operators(tenant_id);
CREATE INDEX idx_operators_shop_floor_id ON operators(shop_floor_id);
CREATE INDEX idx_operators_customer_id ON operators(customer_id);
CREATE INDEX idx_operators_code ON operators(code);
CREATE INDEX idx_operators_is_active ON operators(is_active);
