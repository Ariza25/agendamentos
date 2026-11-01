CREATE TABLE IF NOT EXISTS catalog.jobs (
  id TEXT PRIMARY KEY,
  tenant_id TEXT NOT NULL DEFAULT '',
  type TEXT NOT NULL,
  payload JSONB NOT NULL DEFAULT '{}',
  status TEXT NOT NULL DEFAULT 'pending',
  attempts INTEGER NOT NULL DEFAULT 0,
  max_attempts INTEGER NOT NULL DEFAULT 5,
  run_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  last_error TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_catalog_jobs_poll
  ON catalog.jobs(status, run_at)
  WHERE status IN ('pending', 'failed');

CREATE TABLE IF NOT EXISTS catalog.calendar_connections (
  id TEXT PRIMARY KEY,
  tenant_id TEXT NOT NULL REFERENCES catalog.tenants(id) ON DELETE CASCADE,
  professional_id TEXT NOT NULL DEFAULT '',
  provider TEXT NOT NULL DEFAULT 'google',
  external_calendar_id TEXT NOT NULL DEFAULT 'primary',
  sync_token TEXT NOT NULL DEFAULT '',
  status TEXT NOT NULL DEFAULT 'active',
  last_synced_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_calendar_connections_tenant
  ON catalog.calendar_connections(tenant_id, professional_id);

