-- Insert default admin user
-- Note: This requires a customer to exist first. You may need to create a default customer.
-- Password: Almogavers#2020 (hashed with bcrypt)

-- First, create a default customer if it doesn't exist
INSERT INTO customers (
    id,
    name,
    email,
    status,
    plan,
    created_at,
    updated_at
) VALUES (
    '00000000-0000-0000-0000-000000000001',
    'System',
    'system@turniq.com',
    'active',
    'enterprise',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
) ON CONFLICT (id) DO NOTHING;

-- Insert default admin user
INSERT INTO users (
    id,
    username,
    email,
    password,
    customer_id,
    is_admin,
    is_active,
    created_at,
    updated_at
) VALUES (
    '00000000-0000-0000-0000-000000000001',
    'admin',
    'admin@turniq.com',
    '$2a$10$gxHVnguwvhLpMFaFZ6WEs.fhdvv5NLvxXd0sbcmqUr1GIUY.gsYW6',
    '00000000-0000-0000-0000-000000000001',
    true,
    true,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
) ON CONFLICT (username) DO NOTHING;
