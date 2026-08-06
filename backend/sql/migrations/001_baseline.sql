CREATE TABLE IF NOT EXISTS catalog.tenants (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  slug TEXT NOT NULL UNIQUE,
  active BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS catalog.users (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  email TEXT NOT NULL,
  phone TEXT NOT NULL DEFAULT '',
  password_hash TEXT NOT NULL,
  role TEXT NOT NULL DEFAULT 'receptionist',
  tenant_id TEXT NOT NULL DEFAULT '',
  auth_version INTEGER NOT NULL DEFAULT 1,
  mfa_enabled BOOLEAN NOT NULL DEFAULT FALSE,
  mfa_secret_encrypted TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_catalog_users_tenant_id ON catalog.users(tenant_id);
CREATE UNIQUE INDEX IF NOT EXISTS users_staff_email_unique
  ON catalog.users (LOWER(email))
  WHERE role IN ('owner', 'admin', 'receptionist', 'professional');

CREATE TABLE IF NOT EXISTS catalog.auth_sessions (
  id TEXT PRIMARY KEY,
  user_id TEXT NOT NULL REFERENCES catalog.users(id) ON DELETE CASCADE,
  client_id TEXT NOT NULL DEFAULT '',
  device_name TEXT NOT NULL DEFAULT '',
  ip_address TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  last_seen_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  expires_at TIMESTAMPTZ NOT NULL,
  revoked_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_auth_sessions_active
  ON catalog.auth_sessions(user_id, last_seen_at DESC)
  WHERE revoked_at IS NULL;

CREATE TABLE IF NOT EXISTS catalog.password_reset_tokens (
  token_hash TEXT PRIMARY KEY,
  challenge_key TEXT NOT NULL DEFAULT '',
  user_id TEXT NOT NULL REFERENCES catalog.users(id) ON DELETE CASCADE,
  expires_at TIMESTAMPTZ NOT NULL,
  used_at TIMESTAMPTZ,
  attempt_count INTEGER NOT NULL DEFAULT 0,
  max_attempts INTEGER NOT NULL DEFAULT 0,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_password_reset_tokens_challenge_active
  ON catalog.password_reset_tokens(challenge_key, created_at DESC)
  WHERE used_at IS NULL AND challenge_key <> '';

CREATE TABLE IF NOT EXISTS catalog.auth_rate_limits (
  key_hash TEXT PRIMARY KEY,
  window_started_at TIMESTAMPTZ NOT NULL,
  request_count INTEGER NOT NULL,
  expires_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_auth_rate_limits_expires_at
  ON catalog.auth_rate_limits(expires_at);

CREATE TABLE IF NOT EXISTS catalog.professionals (
  id TEXT PRIMARY KEY,
  tenant_id TEXT NOT NULL REFERENCES catalog.tenants(id) ON DELETE CASCADE,
  user_id TEXT NOT NULL DEFAULT '',
  name TEXT NOT NULL,
  email TEXT NOT NULL DEFAULT '',
  phone TEXT NOT NULL DEFAULT '',
  color TEXT NOT NULL DEFAULT '',
  active BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_catalog_professionals_tenant
  ON catalog.professionals(tenant_id);

CREATE TABLE IF NOT EXISTS catalog.services (
  id TEXT PRIMARY KEY,
  tenant_id TEXT NOT NULL REFERENCES catalog.tenants(id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  description TEXT NOT NULL DEFAULT '',
  duration_minutes INTEGER NOT NULL DEFAULT 30,
  price_cents INTEGER NOT NULL DEFAULT 0,
  color TEXT NOT NULL DEFAULT '',
  active BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_catalog_services_tenant
  ON catalog.services(tenant_id);

CREATE TABLE IF NOT EXISTS catalog.customers (
  id TEXT PRIMARY KEY,
  tenant_id TEXT NOT NULL REFERENCES catalog.tenants(id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  phone TEXT NOT NULL DEFAULT '',
  email TEXT NOT NULL DEFAULT '',
  birthday TEXT NOT NULL DEFAULT '',
  notes TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_catalog_customers_tenant
  ON catalog.customers(tenant_id);

CREATE TABLE IF NOT EXISTS catalog.appointments (
  id TEXT PRIMARY KEY,
  tenant_id TEXT NOT NULL REFERENCES catalog.tenants(id) ON DELETE CASCADE,
  professional_id TEXT NOT NULL,
  service_id TEXT NOT NULL DEFAULT '',
  customer_id TEXT NOT NULL,
  starts_at TIMESTAMPTZ NOT NULL,
  ends_at TIMESTAMPTZ NOT NULL,
  status TEXT NOT NULL DEFAULT 'scheduled',
  notes TEXT NOT NULL DEFAULT '',
  created_by TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_catalog_appointments_tenant_range
  ON catalog.appointments(tenant_id, starts_at, ends_at);

CREATE INDEX IF NOT EXISTS idx_catalog_appointments_professional_range
  ON catalog.appointments(professional_id, starts_at, ends_at);

CREATE INDEX IF NOT EXISTS idx_catalog_appointments_customer_range
  ON catalog.appointments(customer_id, starts_at);

CREATE TABLE IF NOT EXISTS catalog.configurations (
  tenant_id TEXT PRIMARY KEY REFERENCES catalog.tenants(id) ON DELETE CASCADE,
  establishment_name TEXT NOT NULL DEFAULT '',
  working_hours JSONB NOT NULL DEFAULT '[]',
  appointment_interval_minutes INTEGER NOT NULL DEFAULT 30,
  default_duration_minutes INTEGER NOT NULL DEFAULT 30,
  allow_online_booking BOOLEAN NOT NULL DEFAULT TRUE,
  advance_notice_hours INTEGER NOT NULL DEFAULT 0,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

INSERT INTO catalog.tenants (id, name, slug, active)
VALUES ('default', 'Estabelecimento Padrao', 'default', TRUE)
ON CONFLICT (id) DO NOTHING;
