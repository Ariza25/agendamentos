CREATE TABLE IF NOT EXISTS catalog.products (
  id TEXT PRIMARY KEY,
  tenant_id TEXT NOT NULL REFERENCES catalog.tenants(id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  sku TEXT NOT NULL DEFAULT '',
  category TEXT NOT NULL DEFAULT 'general',
  price_cents INTEGER NOT NULL DEFAULT 0,
  cost_cents INTEGER NOT NULL DEFAULT 0,
  stock_quantity NUMERIC(12,2) NOT NULL DEFAULT 0,
  min_stock_alert NUMERIC(12,2) NOT NULL DEFAULT 0,
  is_supply BOOLEAN NOT NULL DEFAULT FALSE,
  unit TEXT NOT NULL DEFAULT 'un',
  active BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_catalog_products_tenant
  ON catalog.products(tenant_id, active);

CREATE TABLE IF NOT EXISTS catalog.service_supplies (
  id TEXT PRIMARY KEY,
  tenant_id TEXT NOT NULL REFERENCES catalog.tenants(id) ON DELETE CASCADE,
  service_id TEXT NOT NULL REFERENCES catalog.services(id) ON DELETE CASCADE,
  product_id TEXT NOT NULL REFERENCES catalog.products(id) ON DELETE CASCADE,
  quantity_consumed NUMERIC(12,2) NOT NULL DEFAULT 1,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_catalog_service_supplies_service
  ON catalog.service_supplies(tenant_id, service_id);

CREATE TABLE IF NOT EXISTS catalog.stock_movements (
  id TEXT PRIMARY KEY,
  tenant_id TEXT NOT NULL REFERENCES catalog.tenants(id) ON DELETE CASCADE,
  product_id TEXT NOT NULL REFERENCES catalog.products(id) ON DELETE CASCADE,
  type TEXT NOT NULL,
  quantity NUMERIC(12,2) NOT NULL,
  unit_cost_cents INTEGER NOT NULL DEFAULT 0,
  notes TEXT NOT NULL DEFAULT '',
  appointment_id TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_catalog_stock_movements_product
  ON catalog.stock_movements(tenant_id, product_id, created_at DESC);

CREATE TABLE IF NOT EXISTS catalog.commission_rules (
  id TEXT PRIMARY KEY,
  tenant_id TEXT NOT NULL REFERENCES catalog.tenants(id) ON DELETE CASCADE,
  professional_id TEXT NOT NULL DEFAULT '',
  service_id TEXT NOT NULL DEFAULT '',
  percentage NUMERIC(5,2) NOT NULL DEFAULT 0,
  fixed_amount_cents INTEGER NOT NULL DEFAULT 0,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_catalog_commission_rules_tenant
  ON catalog.commission_rules(tenant_id);

CREATE TABLE IF NOT EXISTS catalog.commission_records (
  id TEXT PRIMARY KEY,
  tenant_id TEXT NOT NULL REFERENCES catalog.tenants(id) ON DELETE CASCADE,
  appointment_id TEXT NOT NULL DEFAULT '',
  professional_id TEXT NOT NULL REFERENCES catalog.professionals(id) ON DELETE CASCADE,
  service_id TEXT NOT NULL DEFAULT '',
  base_amount_cents INTEGER NOT NULL DEFAULT 0,
  commission_amount_cents INTEGER NOT NULL DEFAULT 0,
  status TEXT NOT NULL DEFAULT 'pending',
  settled_at TIMESTAMPTZ,
  notes TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_catalog_commission_records_prof
  ON catalog.commission_records(tenant_id, professional_id, status);

CREATE TABLE IF NOT EXISTS catalog.financial_transactions (
  id TEXT PRIMARY KEY,
  tenant_id TEXT NOT NULL REFERENCES catalog.tenants(id) ON DELETE CASCADE,
  type TEXT NOT NULL,
  category TEXT NOT NULL DEFAULT 'other',
  amount_cents INTEGER NOT NULL DEFAULT 0,
  description TEXT NOT NULL,
  reference_id TEXT NOT NULL DEFAULT '',
  due_date DATE NOT NULL DEFAULT CURRENT_DATE,
  paid_at TIMESTAMPTZ,
  status TEXT NOT NULL DEFAULT 'pending',
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_catalog_financial_transactions_range
  ON catalog.financial_transactions(tenant_id, due_date, status);

