-- Create jobs table
CREATE TABLE IF NOT EXISTS jobs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    shop_floor_id UUID NOT NULL REFERENCES shopfloors(id) ON DELETE CASCADE,
    workcenter_id UUID NOT NULL,
    job_code VARCHAR(100) NOT NULL,
    product_code VARCHAR(100),
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes
CREATE INDEX idx_jobs_tenant_id ON jobs(tenant_id);
CREATE INDEX idx_jobs_shop_floor_id ON jobs(shop_floor_id);
CREATE INDEX idx_jobs_workcenter_id ON jobs(workcenter_id);
CREATE INDEX idx_jobs_job_code ON jobs(job_code);
