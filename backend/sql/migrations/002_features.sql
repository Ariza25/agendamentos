ALTER TABLE catalog.configurations
  ADD COLUMN IF NOT EXISTS buffer_minutes INTEGER NOT NULL DEFAULT 0,
  ADD COLUMN IF NOT EXISTS closed_days JSONB NOT NULL DEFAULT '[]',
  ADD COLUMN IF NOT EXISTS notifications_enabled BOOLEAN NOT NULL DEFAULT FALSE,
  ADD COLUMN IF NOT EXISTS notification_reminder_hours INTEGER NOT NULL DEFAULT 24,
  ADD COLUMN IF NOT EXISTS notification_channel TEXT NOT NULL DEFAULT 'whatsapp',
  ADD COLUMN IF NOT EXISTS deposit_enabled BOOLEAN NOT NULL DEFAULT FALSE,
  ADD COLUMN IF NOT EXISTS deposit_type TEXT NOT NULL DEFAULT 'fixed',
  ADD COLUMN IF NOT EXISTS deposit_amount_cents INTEGER NOT NULL DEFAULT 0,
  ADD COLUMN IF NOT EXISTS deposit_percentage INTEGER NOT NULL DEFAULT 0,
  ADD COLUMN IF NOT EXISTS no_show_fee_cents INTEGER NOT NULL DEFAULT 0,
  ADD COLUMN IF NOT EXISTS cancellation_policy_hours INTEGER NOT NULL DEFAULT 0,
  ADD COLUMN IF NOT EXISTS calendar_sync_enabled BOOLEAN NOT NULL DEFAULT FALSE;

ALTER TABLE catalog.appointments
  ADD COLUMN IF NOT EXISTS management_code_hash TEXT,
  ADD COLUMN IF NOT EXISTS series_id TEXT NOT NULL DEFAULT '',
  ADD COLUMN IF NOT EXISTS check_in_at TIMESTAMPTZ,
  ADD COLUMN IF NOT EXISTS no_show_at TIMESTAMPTZ,
  ADD COLUMN IF NOT EXISTS payment_status TEXT NOT NULL DEFAULT '',
  ADD COLUMN IF NOT EXISTS no_show_fee_cents INTEGER NOT NULL DEFAULT 0,
  ADD COLUMN IF NOT EXISTS reminder_sent_at TIMESTAMPTZ,
  ADD COLUMN IF NOT EXISTS source TEXT NOT NULL DEFAULT '',
  ADD COLUMN IF NOT EXISTS cancellation_reason TEXT NOT NULL DEFAULT '';

CREATE UNIQUE INDEX IF NOT EXISTS idx_catalog_appointments_management_code
  ON catalog.appointments(management_code_hash)
  WHERE management_code_hash IS NOT NULL;

CREATE TABLE IF NOT EXISTS catalog.service_professionals (
  tenant_id TEXT NOT NULL REFERENCES catalog.tenants(id) ON DELETE CASCADE,
  service_id TEXT NOT NULL REFERENCES catalog.services(id) ON DELETE CASCADE,
  professional_id TEXT NOT NULL REFERENCES catalog.professionals(id) ON DELETE CASCADE,
  PRIMARY KEY (service_id, professional_id)
);

CREATE TABLE IF NOT EXISTS catalog.notification_logs (
  id TEXT PRIMARY KEY,
  tenant_id TEXT NOT NULL REFERENCES catalog.tenants(id) ON DELETE CASCADE,
  appointment_id TEXT NOT NULL DEFAULT '',
  type TEXT NOT NULL,
  channel TEXT NOT NULL DEFAULT '',
  recipient TEXT NOT NULL DEFAULT '',
  subject TEXT NOT NULL DEFAULT '',
  body TEXT NOT NULL DEFAULT '',
  status TEXT NOT NULL DEFAULT 'sent',
  error TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  sent_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_notification_logs_tenant
  ON catalog.notification_logs(tenant_id, created_at DESC);

CREATE TABLE IF NOT EXISTS catalog.payments (
  id TEXT PRIMARY KEY,
  tenant_id TEXT NOT NULL REFERENCES catalog.tenants(id) ON DELETE CASCADE,
  appointment_id TEXT NOT NULL DEFAULT '',
  amount_cents INTEGER NOT NULL DEFAULT 0,
  method TEXT NOT NULL DEFAULT '',
  provider TEXT NOT NULL DEFAULT '',
  provider_payment_id TEXT NOT NULL DEFAULT '',
  status TEXT NOT NULL DEFAULT 'pending',
  description TEXT NOT NULL DEFAULT '',
  paid_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_payments_appointment
  ON catalog.payments(appointment_id);

CREATE INDEX IF NOT EXISTS idx_payments_tenant
  ON catalog.payments(tenant_id, created_at DESC);

CREATE TABLE IF NOT EXISTS catalog.waiting_list (
  id TEXT PRIMARY KEY,
  tenant_id TEXT NOT NULL REFERENCES catalog.tenants(id) ON DELETE CASCADE,
  service_id TEXT NOT NULL DEFAULT '',
  professional_id TEXT NOT NULL DEFAULT '',
  customer_id TEXT NOT NULL DEFAULT '',
  name TEXT NOT NULL DEFAULT '',
  phone TEXT NOT NULL DEFAULT '',
  status TEXT NOT NULL DEFAULT 'waiting',
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  notified_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_waiting_list_tenant
  ON catalog.waiting_list(tenant_id, status);

CREATE INDEX IF NOT EXISTS idx_waiting_list_slot
  ON catalog.waiting_list(tenant_id, service_id, professional_id)
  WHERE status = 'waiting';

CREATE TABLE IF NOT EXISTS catalog.packages (
  id TEXT PRIMARY KEY,
  tenant_id TEXT NOT NULL REFERENCES catalog.tenants(id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  description TEXT NOT NULL DEFAULT '',
  price_cents INTEGER NOT NULL DEFAULT 0,
  sessions_count INTEGER NOT NULL DEFAULT 1,
  duration_days INTEGER NOT NULL DEFAULT 0,
  active BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_packages_tenant
  ON catalog.packages(tenant_id);

CREATE TABLE IF NOT EXISTS catalog.package_credits (
  id TEXT PRIMARY KEY,
  tenant_id TEXT NOT NULL REFERENCES catalog.tenants(id) ON DELETE CASCADE,
  package_id TEXT NOT NULL DEFAULT '',
  customer_id TEXT NOT NULL DEFAULT '',
  sessions_total INTEGER NOT NULL DEFAULT 0,
  sessions_used INTEGER NOT NULL DEFAULT 0,
  expires_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_package_credits_customer
  ON catalog.package_credits(customer_id);

CREATE TABLE IF NOT EXISTS catalog.gift_cards (
  id TEXT PRIMARY KEY,
  tenant_id TEXT NOT NULL REFERENCES catalog.tenants(id) ON DELETE CASCADE,
  code TEXT NOT NULL UNIQUE,
  value_cents INTEGER NOT NULL DEFAULT 0,
  balance_cents INTEGER NOT NULL DEFAULT 0,
  customer_id TEXT NOT NULL DEFAULT '',
  status TEXT NOT NULL DEFAULT 'active',
  expires_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_gift_cards_tenant
  ON catalog.gift_cards(tenant_id);

CREATE TABLE IF NOT EXISTS catalog.memberships (
  id TEXT PRIMARY KEY,
  tenant_id TEXT NOT NULL REFERENCES catalog.tenants(id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  description TEXT NOT NULL DEFAULT '',
  price_cents INTEGER NOT NULL DEFAULT 0,
  billing_period_days INTEGER NOT NULL DEFAULT 30,
  active BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_memberships_tenant
  ON catalog.memberships(tenant_id);

CREATE TABLE IF NOT EXISTS catalog.member_plans (
  id TEXT PRIMARY KEY,
  tenant_id TEXT NOT NULL REFERENCES catalog.tenants(id) ON DELETE CASCADE,
  membership_id TEXT NOT NULL DEFAULT '',
  customer_id TEXT NOT NULL DEFAULT '',
  status TEXT NOT NULL DEFAULT 'active',
  started_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  expires_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_member_plans_customer
  ON catalog.member_plans(customer_id);
