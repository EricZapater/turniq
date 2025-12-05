-- Remove tenant_id column from customers table
ALTER TABLE customers DROP COLUMN IF EXISTS tenant_id;
