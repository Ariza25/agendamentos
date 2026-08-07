CREATE TABLE IF NOT EXISTS catalog.reviews (
  id TEXT PRIMARY KEY,
  tenant_id TEXT NOT NULL REFERENCES catalog.tenants(id) ON DELETE CASCADE,
  appointment_id TEXT NOT NULL DEFAULT '',
  customer_id TEXT NOT NULL DEFAULT '',
  status TEXT NOT NULL DEFAULT 'invited',
  rating INTEGER NOT NULL DEFAULT 0,
  comment TEXT NOT NULL DEFAULT '',
  invited_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  submitted_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_reviews_tenant_created
  ON catalog.reviews(tenant_id, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_reviews_tenant_status
  ON catalog.reviews(tenant_id, status);
