-- =====================================================================
-- ERP Agendamentos - Schema Unificado de Produção
-- =====================================================================
-- Schema canônico único contendo todas as tabelas, índices e dados iniciais
-- para execução idempotente (CREATE IF NOT EXISTS / ON CONFLICT DO NOTHING).
-- =====================================================================

CREATE SCHEMA IF NOT EXISTS catalog;

-- ---------------------------------------------------------------------
-- 1. Tenants (Multi-tenant)
-- ---------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS catalog.tenants (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  slug TEXT NOT NULL UNIQUE,
  segment TEXT NOT NULL DEFAULT 'general',
  active BOOLEAN NOT NULL DEFAULT TRUE,
  plan TEXT NOT NULL DEFAULT 'trial',
  subscription_status TEXT NOT NULL DEFAULT 'trialing',
  trial_ends_at TIMESTAMPTZ DEFAULT (NOW() + INTERVAL '14 days'),
  subscription_expires_at TIMESTAMPTZ,
  active_modules JSONB NOT NULL DEFAULT '["agenda", "appointments", "customers", "services", "professionals", "finance", "commissions", "inventory", "pos_cashregister", "notifications", "reviews"]'::jsonb,
  max_professionals INT NOT NULL DEFAULT 5,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMPTZ
);

-- ---------------------------------------------------------------------
-- 2. Usuários & Sessões de Autenticação
-- ---------------------------------------------------------------------
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

-- ---------------------------------------------------------------------
-- 3. Configurações & Branding do Estabelecimento
-- ---------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS catalog.configurations (
  tenant_id TEXT PRIMARY KEY REFERENCES catalog.tenants(id) ON DELETE CASCADE,
  establishment_name TEXT NOT NULL DEFAULT '',
  working_hours JSONB NOT NULL DEFAULT '[]',
  appointment_interval_minutes INTEGER NOT NULL DEFAULT 30,
  default_duration_minutes INTEGER NOT NULL DEFAULT 30,
  allow_online_booking BOOLEAN NOT NULL DEFAULT TRUE,
  advance_notice_hours INTEGER NOT NULL DEFAULT 0,
  buffer_minutes INTEGER NOT NULL DEFAULT 0,
  closed_days JSONB NOT NULL DEFAULT '[]',
  notifications_enabled BOOLEAN NOT NULL DEFAULT FALSE,
  notification_reminder_hours INTEGER NOT NULL DEFAULT 24,
  notification_channel TEXT NOT NULL DEFAULT 'whatsapp',
  deposit_enabled BOOLEAN NOT NULL DEFAULT FALSE,
  deposit_type TEXT NOT NULL DEFAULT 'fixed',
  deposit_amount_cents INTEGER NOT NULL DEFAULT 0,
  deposit_percentage INTEGER NOT NULL DEFAULT 0,
  no_show_fee_cents INTEGER NOT NULL DEFAULT 0,
  cancellation_policy_hours INTEGER NOT NULL DEFAULT 0,
  calendar_sync_enabled BOOLEAN NOT NULL DEFAULT FALSE,
  logo_url TEXT DEFAULT '',
  primary_color VARCHAR(30) DEFAULT '#b42318',
  address VARCHAR(255) DEFAULT '',
  address_number VARCHAR(50) DEFAULT '',
  neighborhood VARCHAR(100) DEFAULT '',
  city VARCHAR(100) DEFAULT '',
  state VARCHAR(50) DEFAULT '',
  postal_code VARCHAR(20) DEFAULT '',
  contact_phone VARCHAR(50) DEFAULT '',
  contact_whatsapp VARCHAR(50) DEFAULT '',
  instagram VARCHAR(100) DEFAULT '',
  facebook VARCHAR(100) DEFAULT '',
  website VARCHAR(255) DEFAULT '',
  description TEXT DEFAULT '',
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ---------------------------------------------------------------------
-- 4. Profissionais, Serviços & Vínculos
-- ---------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS catalog.professionals (
  id TEXT PRIMARY KEY,
  tenant_id TEXT NOT NULL REFERENCES catalog.tenants(id) ON DELETE CASCADE,
  user_id TEXT NOT NULL DEFAULT '',
  name TEXT NOT NULL,
  email TEXT NOT NULL DEFAULT '',
  phone TEXT NOT NULL DEFAULT '',
  color TEXT NOT NULL DEFAULT '',
  commission_percentage NUMERIC(5,2) DEFAULT 0,
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
  max_capacity INTEGER NOT NULL DEFAULT 1,
  active BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_catalog_services_tenant
  ON catalog.services(tenant_id);

CREATE TABLE IF NOT EXISTS catalog.service_professionals (
  tenant_id TEXT NOT NULL REFERENCES catalog.tenants(id) ON DELETE CASCADE,
  service_id TEXT NOT NULL REFERENCES catalog.services(id) ON DELETE CASCADE,
  professional_id TEXT NOT NULL REFERENCES catalog.professionals(id) ON DELETE CASCADE,
  PRIMARY KEY (service_id, professional_id)
);

-- ---------------------------------------------------------------------
-- 5. Clientes & Recursos (Salas/Cadeiras)
-- ---------------------------------------------------------------------
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

-- ---------------------------------------------------------------------
-- 6. Agendamentos
-- ---------------------------------------------------------------------
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
  source TEXT NOT NULL DEFAULT '',
  cancellation_reason TEXT NOT NULL DEFAULT '',
  payment_status TEXT NOT NULL DEFAULT '',
  no_show_fee_cents INTEGER NOT NULL DEFAULT 0,
  management_code_hash TEXT,
  series_id TEXT NOT NULL DEFAULT '',
  check_in_at TIMESTAMPTZ,
  no_show_at TIMESTAMPTZ,
  reminder_sent_at TIMESTAMPTZ,
  resource_id TEXT NOT NULL DEFAULT '',
  recurrence_rule TEXT NOT NULL DEFAULT '',
  timezone TEXT NOT NULL DEFAULT '',
  service_ids JSONB NOT NULL DEFAULT '[]',
  lgpd_accepted_at TIMESTAMPTZ,
  marketing_consent BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_catalog_appointments_tenant_range
  ON catalog.appointments(tenant_id, starts_at, ends_at);

CREATE INDEX IF NOT EXISTS idx_catalog_appointments_professional_range
  ON catalog.appointments(professional_id, starts_at, ends_at);

CREATE INDEX IF NOT EXISTS idx_catalog_appointments_customer_range
  ON catalog.appointments(customer_id, starts_at);

CREATE UNIQUE INDEX IF NOT EXISTS idx_catalog_appointments_management_code
  ON catalog.appointments(management_code_hash)
  WHERE management_code_hash IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_appointments_resource
  ON catalog.appointments(tenant_id, resource_id, starts_at);

-- ---------------------------------------------------------------------
-- 7. Notificações & Mensageria
-- ---------------------------------------------------------------------
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

-- ---------------------------------------------------------------------
-- 8. Pagamentos (Booking & Mercado Pago Pix)
-- ---------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS catalog.payments (
  id TEXT PRIMARY KEY,
  tenant_id TEXT NOT NULL REFERENCES catalog.tenants(id) ON DELETE CASCADE,
  appointment_id TEXT NOT NULL DEFAULT '',
  amount_cents INTEGER NOT NULL DEFAULT 0,
  method TEXT NOT NULL DEFAULT '',
  provider TEXT NOT NULL DEFAULT '',
  provider_payment_id TEXT NOT NULL DEFAULT '',
  pix_copia_e_cola TEXT NOT NULL DEFAULT '',
  pix_qr_code_base64 TEXT NOT NULL DEFAULT '',
  provider_data TEXT NOT NULL DEFAULT '',
  status TEXT NOT NULL DEFAULT 'pending',
  description TEXT NOT NULL DEFAULT '',
  paid_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_payments_appointment
  ON catalog.payments(appointment_id);

CREATE INDEX IF NOT EXISTS idx_payments_tenant
  ON catalog.payments(tenant_id, created_at DESC);

-- ---------------------------------------------------------------------
-- 9. Lista de Espera
-- ---------------------------------------------------------------------
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

-- ---------------------------------------------------------------------
-- 10. Fidelidade, Pacotes, Vouchers & Assinaturas
-- ---------------------------------------------------------------------
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

-- ---------------------------------------------------------------------
-- 11. Avaliações (Reviews)
-- ---------------------------------------------------------------------
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

-- ---------------------------------------------------------------------
-- 12. Estoque & Produtos (Venda & Insumos de Serviços)
-- ---------------------------------------------------------------------
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

-- ---------------------------------------------------------------------
-- 13. Frente de Caixa (Caixa Diário, Sangrias & Suprimentos)
-- ---------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS catalog.cash_registers (
  id TEXT PRIMARY KEY,
  tenant_id TEXT NOT NULL REFERENCES catalog.tenants(id) ON DELETE CASCADE,
  opened_by_user_id TEXT NOT NULL DEFAULT '',
  closed_by_user_id TEXT NOT NULL DEFAULT '',
  status VARCHAR(20) NOT NULL DEFAULT 'open',
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

CREATE TABLE IF NOT EXISTS catalog.cash_register_entries (
  id TEXT PRIMARY KEY,
  tenant_id TEXT NOT NULL REFERENCES catalog.tenants(id) ON DELETE CASCADE,
  cash_register_id TEXT NOT NULL REFERENCES catalog.cash_registers(id) ON DELETE CASCADE,
  type VARCHAR(20) NOT NULL,
  amount_cents INTEGER NOT NULL DEFAULT 0,
  reason TEXT NOT NULL DEFAULT '',
  created_by_user_id TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_catalog_cash_reg_entries_reg
  ON catalog.cash_register_entries(tenant_id, cash_register_id);

-- ---------------------------------------------------------------------
-- 14. Gestão Financeira (Receitas & Despesas)
-- ---------------------------------------------------------------------
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
  cash_register_id TEXT DEFAULT '',
  payment_method VARCHAR(30) DEFAULT 'other',
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_catalog_financial_transactions_range
  ON catalog.financial_transactions(tenant_id, due_date, status);

-- ---------------------------------------------------------------------
-- 15. Comissões & Repasses da Equipe
-- ---------------------------------------------------------------------
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

CREATE TABLE IF NOT EXISTS catalog.commission_records (
  id TEXT PRIMARY KEY,
  tenant_id TEXT NOT NULL REFERENCES catalog.tenants(id) ON DELETE CASCADE,
  appointment_id TEXT NOT NULL DEFAULT '',
  professional_id TEXT NOT NULL REFERENCES catalog.professionals(id) ON DELETE CASCADE,
  service_id TEXT NOT NULL DEFAULT '',
  base_amount_cents INTEGER NOT NULL DEFAULT 0,
  commission_amount_cents INTEGER NOT NULL DEFAULT 0,
  status TEXT NOT NULL DEFAULT 'pending',
  payout_id TEXT REFERENCES catalog.commission_payouts(id) ON DELETE SET NULL,
  settled_at TIMESTAMPTZ,
  notes TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_catalog_commission_records_prof
  ON catalog.commission_records(tenant_id, professional_id, status);

CREATE INDEX IF NOT EXISTS idx_commission_records_payout
  ON catalog.commission_records(tenant_id, payout_id);

-- ---------------------------------------------------------------------
-- 16. Filas de Jobs & Integrações de Calendário Externo
-- ---------------------------------------------------------------------
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

-- ---------------------------------------------------------------------
-- 17. Seed / Carga Inicial Padrão (Idempotente)
-- ---------------------------------------------------------------------
-- Tenant Master/Default
INSERT INTO catalog.tenants (
  id, name, slug, segment, active, plan, subscription_status,
  trial_ends_at, subscription_expires_at, active_modules, max_professionals
) VALUES (
  'default', 'Estabelecimento Padrão', 'default', 'general', TRUE, 'enterprise', 'active',
  NULL, (NOW() + INTERVAL '10 years'),
  '["agenda", "appointments", "customers", "services", "professionals", "finance", "commissions", "inventory", "pos_cashregister", "commerce", "reviews", "notifications", "widget", "calendar_sync"]'::jsonb,
  999
) ON CONFLICT (id) DO UPDATE SET
  plan = EXCLUDED.plan,
  subscription_status = EXCLUDED.subscription_status,
  active_modules = EXCLUDED.active_modules,
  max_professionals = EXCLUDED.max_professionals;

-- Configuração do Tenant Default
INSERT INTO catalog.configurations (
  tenant_id, establishment_name, working_hours, appointment_interval_minutes,
  default_duration_minutes, allow_online_booking, advance_notice_hours
) VALUES (
  'default', 'Estabelecimento Padrão',
  '[{"dayOfWeek":1,"openTime":"09:00","closeTime":"18:00","closed":false},{"dayOfWeek":2,"openTime":"09:00","closeTime":"18:00","closed":false},{"dayOfWeek":3,"openTime":"09:00","closeTime":"18:00","closed":false},{"dayOfWeek":4,"openTime":"09:00","closeTime":"18:00","closed":false},{"dayOfWeek":5,"openTime":"09:00","closeTime":"18:00","closed":false},{"dayOfWeek":6,"openTime":"09:00","closeTime":"14:00","closed":false},{"dayOfWeek":0,"openTime":"00:00","closeTime":"00:00","closed":true}]'::jsonb,
  30, 30, TRUE, 0
) ON CONFLICT (tenant_id) DO NOTHING;

-- Serviços de Demonstração
INSERT INTO catalog.services (id, tenant_id, name, description, duration_minutes, price_cents, color, active)
VALUES
  ('demo-svc-corte', 'default', 'Corte de Cabelo', 'Corte na tesoura ou máquina com finalização.', 30, 5000, '#3b82f6', TRUE),
  ('demo-svc-barba', 'default', 'Barba', 'Aparo e modelagem da barba.', 30, 4000, '#8b5cf6', TRUE),
  ('demo-svc-corte-barba', 'default', 'Corte + Barba', 'Combo corte e barba completo.', 60, 8000, '#6366f1', TRUE),
  ('demo-svc-coloracao', 'default', 'Coloração', 'Coloração completa com lavagem.', 120, 15000, '#ef4444', TRUE),
  ('demo-svc-sobrancelha', 'default', 'Sobrancelha', 'Design de sobrancelha.', 30, 3000, '#f59e0b', TRUE)
ON CONFLICT (id) DO NOTHING;

-- Profissionais de Demonstração
INSERT INTO catalog.professionals (id, tenant_id, user_id, name, email, phone, color, commission_percentage, active)
VALUES
  ('demo-prof-ana', 'default', '', 'Ana Souza', 'ana@demo.local', '(11) 90000-0001', '#3b82f6', 30.00, TRUE),
  ('demo-prof-bruno', 'default', '', 'Bruno Lima', 'bruno@demo.local', '(11) 90000-0002', '#10b981', 35.00, TRUE),
  ('demo-prof-carla', 'default', '', 'Carla Mendes', 'carla@demo.local', '(11) 90000-0003', '#f43f5e', 40.00, TRUE)
ON CONFLICT (id) DO NOTHING;

-- Vínculo Serviços x Profissionais
INSERT INTO catalog.service_professionals (tenant_id, service_id, professional_id)
VALUES
  ('default', 'demo-svc-corte', 'demo-prof-ana'),
  ('default', 'demo-svc-corte', 'demo-prof-bruno'),
  ('default', 'demo-svc-barba', 'demo-prof-bruno'),
  ('default', 'demo-svc-corte-barba', 'demo-prof-bruno'),
  ('default', 'demo-svc-corte-barba', 'demo-prof-carla'),
  ('default', 'demo-svc-coloracao', 'demo-prof-carla'),
  ('default', 'demo-svc-sobrancelha', 'demo-prof-carla')
ON CONFLICT (service_id, professional_id) DO NOTHING;

-- Clientes de Demonstração
INSERT INTO catalog.customers (id, tenant_id, name, phone, email, birthday, notes)
VALUES
  ('demo-cust-joao', 'default', 'João Silva', '(11) 98888-0001', 'joao.silva@demo.local', '1990-03-15', 'Cliente desde o início.'),
  ('demo-cust-maria', 'default', 'Maria Oliveira', '(11) 98888-0002', 'maria.oliveira@demo.local', '1992-07-22', ''),
  ('demo-cust-pedro', 'default', 'Pedro Santos', '(11) 98888-0003', 'pedro.santos@demo.local', '1985-11-02', 'Prefere corte na tesoura.'),
  ('demo-cust-fernanda', 'default', 'Fernanda Costa', '(11) 98888-0004', 'fernanda.costa@demo.local', '1995-01-30', ''),
  ('demo-cust-rafael', 'default', 'Rafael Almeida', '(11) 98888-0005', 'rafael.almeida@demo.local', '1988-09-12', ''),
  ('demo-cust-juliana', 'default', 'Juliana Rocha', '(11) 98888-0006', 'juliana.rocha@demo.local', '1993-05-08', '')
ON CONFLICT (id) DO NOTHING;

