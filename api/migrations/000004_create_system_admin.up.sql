-- Enable pgcrypto extension for password hashing
CREATE EXTENSION IF NOT EXISTS pgcrypto;

DO $$
DECLARE
    v_customer_id UUID;
    v_user_id UUID;
    v_password_plain TEXT := 'RawCraft_2026';
BEGIN
    -- 1. Create System Customer
    INSERT INTO customers (name, email, language, status, max_operators, max_workcenters, max_shop_floors, max_users, max_jobs)
    VALUES ('System', 'admin@turniq.com', 'en', 'active', 999, 999, 999, 999, 999)
    RETURNING id INTO v_customer_id;

    -- 2. Create Admin User (Password: admin1234)
    INSERT INTO users (customer_id, email, password, username, is_admin, is_active)
    VALUES (v_customer_id, 'admin@turniq.com', crypt(v_password_plain, gen_salt('bf')), 'System Admin', true, true)
    RETURNING id INTO v_user_id;

END $$;
