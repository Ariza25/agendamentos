-- 010_branding_pos_and_cash_registers.sql

-- 1. Tenant Branding & Public Page Details in Configurations
ALTER TABLE catalog.configurations ADD COLUMN IF NOT EXISTS logo_url TEXT DEFAULT '';
ALTER TABLE catalog.configurations ADD COLUMN IF NOT EXISTS primary_color VARCHAR(30) DEFAULT '#b42318';
ALTER TABLE catalog.configurations ADD COLUMN IF NOT EXISTS address VARCHAR(255) DEFAULT '';
ALTER TABLE catalog.configurations ADD COLUMN IF NOT EXISTS address_number VARCHAR(50) DEFAULT '';
ALTER TABLE catalog.configurations ADD COLUMN IF NOT EXISTS neighborhood VARCHAR(100) DEFAULT '';
ALTER TABLE catalog.configurations ADD COLUMN IF NOT EXISTS city VARCHAR(100) DEFAULT '';
ALTER TABLE catalog.configurations ADD COLUMN IF NOT EXISTS state VARCHAR(50) DEFAULT '';
ALTER TABLE catalog.configurations ADD COLUMN IF NOT EXISTS postal_code VARCHAR(20) DEFAULT '';
ALTER TABLE catalog.configurations ADD COLUMN IF NOT EXISTS contact_phone VARCHAR(50) DEFAULT '';
ALTER TABLE catalog.configurations ADD COLUMN IF NOT EXISTS contact_whatsapp VARCHAR(50) DEFAULT '';
ALTER TABLE catalog.configurations ADD COLUMN IF NOT EXISTS instagram VARCHAR(100) DEFAULT '';
ALTER TABLE catalog.configurations ADD COLUMN IF NOT EXISTS facebook VARCHAR(100) DEFAULT '';
ALTER TABLE catalog.configurations ADD COLUMN IF NOT EXISTS website VARCHAR(255) DEFAULT '';
ALTER TABLE catalog.configurations ADD COLUMN IF NOT EXISTS description TEXT DEFAULT '';

-- 2. Professional Default Commission
ALTER TABLE catalog.professionals ADD COLUMN IF NOT EXISTS commission_percentage NUMERIC(5,2) DEFAULT 0;

-- 3. Cash Registers (Controle de Abertura / Sangria / Fechamento de Caixa Diário)
CREATE TABLE IF NOT EXISTS catalog.cash_registers (
  id TEXT PRIMARY KEY,
  tenant_id TEXT NOT NULL REFERENCES catalog.tenants(id) ON DELETE CASCADE,
  opened_by_user_id TEXT NOT NULL DEFAULT '',
  closed_by_user_id TEXT NOT NULL DEFAULT '',
  status VARCHAR(20) NOT NULL DEFAULT 'open', -- 'open', 'closed'
  initial_balance_cents INTEGER NOT NULL DEFAULT 0,
  closing_balance_cents INTEGER NOT NULL DEFAULT 0,
  system_balance_cents INTEGER NOT NULL DEFAULT 0,
  difference_cents INTEGER NOT NULL DEFAULT 0,
  opened_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  closed_at TIMESTAMPTZ,
  notes TEXT NOT NULL DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_catalog_cash_registers_tenant
  ON catalog.cash_registers(tenant_id, status, opened_at DESC);

-- Cash Register Entries (Sangrias e Suprimentos)
CREATE TABLE IF NOT EXISTS catalog.cash_register_entries (
  id TEXT PRIMARY KEY,
  tenant_id TEXT NOT NULL REFERENCES catalog.tenants(id) ON DELETE CASCADE,
  cash_register_id TEXT NOT NULL REFERENCES catalog.cash_registers(id) ON DELETE CASCADE,
  type VARCHAR(20) NOT NULL, -- 'sangria' (retirada) ou 'suprimento' (aporte)
  amount_cents INTEGER NOT NULL DEFAULT 0,
  reason TEXT NOT NULL DEFAULT '',
  created_by_user_id TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_catalog_cash_reg_entries_reg
  ON catalog.cash_register_entries(tenant_id, cash_register_id);

-- 4. Financial Transactions Link to Cash Register and Payment Method
ALTER TABLE catalog.financial_transactions ADD COLUMN IF NOT EXISTS cash_register_id TEXT DEFAULT '';
ALTER TABLE catalog.financial_transactions ADD COLUMN IF NOT EXISTS payment_method VARCHAR(30) DEFAULT 'other';

