INSERT INTO catalog.services (id, tenant_id, name, description, duration_minutes, price_cents, color, active)
VALUES
  ('demo-svc-corte', 'default', 'Corte de Cabelo', 'Corte na tesoura ou máquina com finalizacao.', 30, 5000, '#3b82f6', TRUE),
  ('demo-svc-barba', 'default', 'Barba', 'Aparo e modelagem da barba.', 30, 4000, '#8b5cf6', TRUE),
  ('demo-svc-corte-barba', 'default', 'Corte + Barba', 'Combo corte e barba completo.', 60, 8000, '#6366f1', TRUE),
  ('demo-svc-coloracao', 'default', 'Coloracao', 'Coloração completa com lavagem.', 120, 15000, '#ef4444', TRUE),
  ('demo-svc-sobrancelha', 'default', 'Sobrancelha', 'Design de sobrancelha.', 30, 3000, '#f59e0b', TRUE)
ON CONFLICT (id) DO NOTHING;

INSERT INTO catalog.professionals (id, tenant_id, user_id, name, email, phone, color, active)
VALUES
  ('demo-prof-ana', 'default', '', 'Ana Souza', 'ana@demo.local', '(11) 90000-0001', '#3b82f6', TRUE),
  ('demo-prof-bruno', 'default', '', 'Bruno Lima', 'bruno@demo.local', '(11) 90000-0002', '#10b981', TRUE),
  ('demo-prof-carla', 'default', '', 'Carla Mendes', 'carla@demo.local', '(11) 90000-0003', '#f43f5e', TRUE)
ON CONFLICT (id) DO NOTHING;

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

INSERT INTO catalog.customers (id, tenant_id, name, phone, email, birthday, notes)
VALUES
  ('demo-cust-joao', 'default', 'João Silva', '(11) 98888-0001', 'joao.silva@demo.local', '1990-03-15', 'Cliente desde o comeco.'),
  ('demo-cust-maria', 'default', 'Maria Oliveira', '(11) 98888-0002', 'maria.oliveira@demo.local', '1992-07-22', ''),
  ('demo-cust-pedro', 'default', 'Pedro Santos', '(11) 98888-0003', 'pedro.santos@demo.local', '1985-11-02', 'Prefere corte na tesoura.'),
  ('demo-cust-fernanda', 'default', 'Fernanda Costa', '(11) 98888-0004', 'fernanda.costa@demo.local', '1995-01-30', ''),
  ('demo-cust-rafael', 'default', 'Rafael Almeida', '(11) 98888-0005', 'rafael.almeida@demo.local', '1988-09-12', ''),
  ('demo-cust-juliana', 'default', 'Juliana Rocha', '(11) 98888-0006', 'juliana.rocha@demo.local', '1993-05-08', '')
ON CONFLICT (id) DO NOTHING;

INSERT INTO catalog.appointments (id, tenant_id, professional_id, service_id, customer_id, starts_at, ends_at, status, notes, created_by, source, payment_status)
VALUES
  ('demo-appt-1', 'default', 'demo-prof-ana', 'demo-svc-corte', 'demo-cust-joao', (CURRENT_DATE - 2) + TIME '09:00', (CURRENT_DATE - 2) + TIME '09:30', 'completed', '', 'admin', 'admin', 'paid'),
  ('demo-appt-2', 'default', 'demo-prof-bruno', 'demo-svc-barba', 'demo-cust-maria', (CURRENT_DATE - 1) + TIME '09:30', (CURRENT_DATE - 1) + TIME '10:00', 'no_show', 'Nao compareceu.', 'admin', 'admin', ''),
  ('demo-appt-3', 'default', 'demo-prof-carla', 'demo-svc-corte-barba', 'demo-cust-pedro', (CURRENT_DATE - 1) + TIME '11:00', (CURRENT_DATE - 1) + TIME '12:00', 'cancelled', 'Cliente cancelou.', 'admin', 'admin', ''),
  ('demo-appt-4', 'default', 'demo-prof-ana', 'demo-svc-corte', 'demo-cust-joao', (CURRENT_DATE - 1) + TIME '14:00', (CURRENT_DATE - 1) + TIME '14:30', 'completed', '', 'admin', 'admin', 'paid'),
  ('demo-appt-5', 'default', 'demo-prof-ana', 'demo-svc-corte', 'demo-cust-maria', CURRENT_DATE + TIME '09:00', CURRENT_DATE + TIME '09:30', 'in_progress', '', 'admin', 'admin', 'paid'),
  ('demo-appt-6', 'default', 'demo-prof-bruno', 'demo-svc-corte-barba', 'demo-cust-pedro', CURRENT_DATE + TIME '10:00', CURRENT_DATE + TIME '11:00', 'confirmed', '', 'admin', 'admin', ''),
  ('demo-appt-7', 'default', 'demo-prof-carla', 'demo-svc-coloracao', 'demo-cust-fernanda', CURRENT_DATE + TIME '11:00', CURRENT_DATE + TIME '13:00', 'scheduled', '', 'admin', 'admin', ''),
  ('demo-appt-8', 'default', 'demo-prof-carla', 'demo-svc-sobrancelha', 'demo-cust-juliana', (CURRENT_DATE + 1) + TIME '09:00', (CURRENT_DATE + 1) + TIME '09:30', 'confirmed', '', 'admin', 'admin', ''),
  ('demo-appt-9', 'default', 'demo-prof-ana', 'demo-svc-corte', 'demo-cust-rafael', (CURRENT_DATE + 1) + TIME '10:00', (CURRENT_DATE + 1) + TIME '10:30', 'scheduled', '', 'booking', 'online', ''),
  ('demo-appt-10', 'default', 'demo-prof-bruno', 'demo-svc-barba', 'demo-cust-juliana', (CURRENT_DATE + 2) + TIME '15:00', (CURRENT_DATE + 2) + TIME '15:30', 'scheduled', '', 'admin', 'admin', '');

INSERT INTO catalog.payments (id, tenant_id, appointment_id, amount_cents, method, provider, provider_payment_id, status, description, paid_at)
VALUES
  ('demo-pay-1', 'default', 'demo-appt-1', 5000, 'pix', 'mock', '', 'paid', 'Corte de Cabelo', (CURRENT_DATE - 2) + TIME '09:05'),
  ('demo-pay-2', 'default', 'demo-appt-5', 5000, 'pix', 'mock', '', 'paid', 'Corte de Cabelo', CURRENT_DATE + TIME '09:02'),
  ('demo-pay-3', 'default', 'demo-appt-6', 8000, 'credit_card', 'mock', '', 'pending', 'Corte + Barba', NULL);

INSERT INTO catalog.waiting_list (id, tenant_id, service_id, professional_id, customer_id, name, phone, status)
VALUES
  ('demo-wl-1', 'default', 'demo-svc-corte', 'demo-prof-ana', 'demo-cust-rafael', 'Rafael Almeida', '(11) 98888-0005', 'waiting');

INSERT INTO catalog.packages (id, tenant_id, name, description, price_cents, sessions_count, duration_days, active)
VALUES
  ('demo-pkg-corte', 'default', 'Pacote Corte Mensal', '4 cortes com validade de 60 dias.', 18000, 4, 60, TRUE),
  ('demo-pkg-barba', 'default', 'Pacote Barba Mensal', '4 barbas com validade de 60 dias.', 14000, 4, 60, TRUE)
ON CONFLICT (id) DO NOTHING;

INSERT INTO catalog.package_credits (id, tenant_id, package_id, customer_id, sessions_total, sessions_used, expires_at)
VALUES
  ('demo-credit-1', 'default', 'demo-pkg-corte', 'demo-cust-joao', 4, 1, (CURRENT_DATE + 60) + TIME '23:59'),
  ('demo-credit-2', 'default', 'demo-pkg-barba', 'demo-cust-pedro', 4, 0, (CURRENT_DATE + 60) + TIME '23:59');

INSERT INTO catalog.gift_cards (id, tenant_id, code, value_cents, balance_cents, customer_id, status, expires_at)
VALUES
  ('demo-gc-1', 'default', 'DEMO-2026-0001', 20000, 15000, 'demo-cust-maria', 'active', (CURRENT_DATE + 180) + TIME '23:59'),
  ('demo-gc-2', 'default', 'DEMO-2026-0002', 10000, 10000, '', 'active', (CURRENT_DATE + 180) + TIME '23:59');

INSERT INTO catalog.memberships (id, tenant_id, name, description, price_cents, billing_period_days, active)
VALUES
  ('demo-membership-1', 'default', 'Assinatura Premium', 'Beneficios mensais com prioridade de horario.', 12000, 30, TRUE)
ON CONFLICT (id) DO NOTHING;

INSERT INTO catalog.member_plans (id, tenant_id, membership_id, customer_id, status, started_at, expires_at)
VALUES
  ('demo-plan-1', 'default', 'demo-membership-1', 'demo-cust-fernanda', 'active', CURRENT_DATE + TIME '00:00', (CURRENT_DATE + 30) + TIME '23:59');

INSERT INTO catalog.notification_logs (id, tenant_id, appointment_id, type, channel, recipient, subject, body, status, sent_at)
VALUES
  ('demo-notif-1', 'default', 'demo-appt-8', 'reminder', 'whatsapp', '(11) 98888-0006', 'Lembrete de agendamento', 'Seu horario e amanha as 09:00.', 'sent', NOW());

INSERT INTO catalog.configurations (tenant_id, establishment_name, working_hours, appointment_interval_minutes, default_duration_minutes, allow_online_booking, advance_notice_hours, buffer_minutes, closed_days, notifications_enabled, notification_reminder_hours, notification_channel, deposit_enabled, deposit_type, deposit_amount_cents, deposit_percentage, no_show_fee_cents, cancellation_policy_hours, calendar_sync_enabled)
VALUES (
  'default',
  'Barbearia Demo',
  '[{"dayOfWeek":0,"openTime":"09:00","closeTime":"18:00","enabled":false},{"dayOfWeek":1,"openTime":"09:00","closeTime":"18:00","enabled":true},{"dayOfWeek":2,"openTime":"09:00","closeTime":"18:00","enabled":true},{"dayOfWeek":3,"openTime":"09:00","closeTime":"18:00","enabled":true},{"dayOfWeek":4,"openTime":"09:00","closeTime":"18:00","enabled":true},{"dayOfWeek":5,"openTime":"09:00","closeTime":"18:00","enabled":true},{"dayOfWeek":6,"openTime":"09:00","closeTime":"18:00","enabled":true}]',
  30,
  30,
  TRUE,
  2,
  10,
  '[]',
  TRUE,
  24,
  'whatsapp',
  FALSE,
  'percentage',
  0,
  0,
  5000,
  12,
  TRUE
)
ON CONFLICT (tenant_id) DO NOTHING;
