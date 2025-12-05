-- Remove tenant_id column from users table
ALTER TABLE users DROP COLUMN IF EXISTS tenant_id;
