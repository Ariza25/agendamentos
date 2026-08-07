CREATE TABLE IF NOT EXISTS catalog.resources (
  id TEXT PRIMARY KEY,
  tenant_id TEXT NOT NULL REFERENCES catalog.tenants(id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  type TEXT NOT NULL DEFAULT 'room',
  capacity INTEGER NOT NULL DEFAULT 1,
  active BOOLEAN NOT NULL DEFAULT TRUE,
  notes TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_resources_tenant
  ON catalog.resources(tenant_id);

ALTER TABLE catalog.appointments
  ADD COLUMN IF NOT EXISTS resource_id TEXT NOT NULL DEFAULT '';

CREATE INDEX IF NOT EXISTS idx_appointments_resource
  ON catalog.appointments(tenant_id, resource_id, starts_at);
