-- Enable pgcrypto extension for password hashing
CREATE EXTENSION IF NOT EXISTS pgcrypto;

WITH new_customer AS (
    INSERT INTO customers (name, email, language, status, max_operators, max_workcenters, max_shop_floors, max_users, max_jobs)
    VALUES ('System', 'admin@turniq.com', 'en', 'active', 999, 999, 999, 999, 999)
    RETURNING id
)
INSERT INTO users (customer_id, email, password, username, is_admin, is_active)
SELECT 
    id, 
    'admin@turniq.com', 
    crypt('RawCraft_2026', gen_salt('bf')), 
    'System Admin', 
    true, 
    true
FROM new_customer;
