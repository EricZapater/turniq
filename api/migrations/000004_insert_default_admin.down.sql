-- Remove default admin user
DELETE FROM users WHERE username = 'admin' AND id = '00000000-0000-0000-0000-000000000001';

-- Remove default system customer (only if no other users reference it)
DELETE FROM customers 
WHERE id = '00000000-0000-0000-0000-000000000001' 
AND NOT EXISTS (
    SELECT 1 FROM users WHERE customer_id = '00000000-0000-0000-0000-000000000001' AND id != '00000000-0000-0000-0000-000000000001'
);
