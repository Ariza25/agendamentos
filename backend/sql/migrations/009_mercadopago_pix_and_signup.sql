-- Migration 009: Mercado Pago Pix, LGPD Consent & Onboarding Segments

ALTER TABLE catalog.payments
  ADD COLUMN IF NOT EXISTS pix_copia_e_cola TEXT NOT NULL DEFAULT '',
  ADD COLUMN IF NOT EXISTS pix_qr_code_base64 TEXT NOT NULL DEFAULT '',
  ADD COLUMN IF NOT EXISTS provider_data TEXT NOT NULL DEFAULT '';

ALTER TABLE catalog.appointments
  ADD COLUMN IF NOT EXISTS lgpd_accepted_at TIMESTAMPTZ,
  ADD COLUMN IF NOT EXISTS marketing_consent BOOLEAN NOT NULL DEFAULT true;

ALTER TABLE catalog.tenants
  ADD COLUMN IF NOT EXISTS segment TEXT NOT NULL DEFAULT 'general';

