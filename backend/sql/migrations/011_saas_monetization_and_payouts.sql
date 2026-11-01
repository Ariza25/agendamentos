-- =====================================================================
-- Migration 011: SaaS Monetization (Trial 14d + Plans + Modules) & Commission Payouts
-- =====================================================================

-- 1. Tenants SaaS Subscription & Module Flags
ALTER TABLE catalog.tenants
  ADD COLUMN IF NOT EXISTS plan TEXT NOT NULL DEFAULT 'trial',
  ADD COLUMN IF NOT EXISTS subscription_status TEXT NOT NULL DEFAULT 'trialing',
  ADD COLUMN IF NOT EXISTS trial_ends_at TIMESTAMPTZ DEFAULT (NOW() + INTERVAL '14 days'),
  ADD COLUMN IF NOT EXISTS subscription_expires_at TIMESTAMPTZ,
  ADD COLUMN IF NOT EXISTS active_modules JSONB NOT NULL DEFAULT '["agenda", "appointments", "customers", "services", "professionals", "finance", "commissions", "inventory", "pos_cashregister", "notifications", "reviews"]'::jsonb,
  ADD COLUMN IF NOT EXISTS max_professionals INT NOT NULL DEFAULT 5;

-- Ensure default / master tenant has enterprise plan and all modules unlocked
UPDATE catalog.tenants
SET
  plan = 'enterprise',
  subscription_status = 'active',
  trial_ends_at = NULL,
  subscription_expires_at = (NOW() + INTERVAL '10 years'),
  active_modules = '["agenda", "appointments", "customers", "services", "professionals", "finance", "commissions", "inventory", "pos_cashregister", "commerce", "reviews", "notifications", "widget", "calendar_sync"]'::jsonb,
  max_professionals = 999
WHERE id = 'default' OR slug = 'default';

-- 2. Commission Payouts (Fechamento / Repasse de Comissões)
CREATE TABLE IF NOT EXISTS catalog.commission_payouts (
  id TEXT PRIMARY KEY,
  tenant_id TEXT NOT NULL REFERENCES catalog.tenants(id) ON DELETE CASCADE,
  professional_id TEXT NOT NULL REFERENCES catalog.professionals(id) ON DELETE CASCADE,
  total_cents INT NOT NULL CHECK (total_cents >= 0),
  records_count INT NOT NULL DEFAULT 0,
  from_date TIMESTAMPTZ,
  to_date TIMESTAMPTZ,
  payment_method TEXT NOT NULL DEFAULT 'pix',
  notes TEXT,
  financial_transaction_id TEXT REFERENCES catalog.financial_transactions(id) ON DELETE SET NULL,
  created_by_user_id TEXT,
  paid_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_commission_payouts_tenant_prof
  ON catalog.commission_payouts(tenant_id, professional_id);

-- 3. Link Commission Records to Payout
ALTER TABLE catalog.commission_records
  ADD COLUMN IF NOT EXISTS payout_id TEXT REFERENCES catalog.commission_payouts(id) ON DELETE SET NULL;

CREATE INDEX IF NOT EXISTS idx_commission_records_payout
  ON catalog.commission_records(tenant_id, payout_id);
