#!/usr/bin/env bash
# =====================================================================
# Script de Emissão e Renovação de Certificados SSL (Let's Encrypt / Certbot)
# =====================================================================
# Automatiza a emissão de certificados SSL gratuitos da Let's Encrypt
# via Certbot em modo standalone ou webroot com Nginx.
#
# Uso:
#   ./scripts/certbot_ssl.sh meu-dominio.com.br admin.meu-dominio.com.br contato@meu-dominio.com.br
# =====================================================================

set -eo pipefail

DOMAIN="${1:-}"
ADMIN_DOMAIN="${2:-}"
EMAIL="${3:-}"

if [[ -z "${DOMAIN}" || -z "${EMAIL}" ]]; then
  echo "Uso: $0 <dominio_principal> [subdominio_admin] <email_responsavel>"
  echo "Exemplo: $0 agendaplus.com.br app.agendaplus.com.br contato@agendaplus.com.br"
  exit 1
fi

DOMAINS_ARGS="-d ${DOMAIN}"
if [[ -n "${ADMIN_DOMAIN}" && "${ADMIN_DOMAIN}" != "${DOMAIN}" ]]; then
  DOMAINS_ARGS="${DOMAINS_ARGS} -d ${ADMIN_DOMAIN}"
fi

echo "=========================================================="
echo "Solicitando certificado SSL para: ${DOMAINS_ARGS}"
echo "Email de aviso/renovação: ${EMAIL}"
echo "=========================================================="

if ! command -v certbot &> /dev/null; then
  echo "[INFO] Certbot não encontrado. Instalando..."
  if command -v apt-get &> /dev/null; then
    sudo apt-get update && sudo apt-get install -y certbot python3-certbot-nginx
  else
    echo "[ERRO] Instale o certbot na sua distribuição ou utilize a imagem docker certbot/certbot."
    exit 1
  fi
fi

# Emite o certificado
certbot certonly \
  --standalone \
  --agree-tos \
  --non-interactive \
  --email "${EMAIL}" \
  ${DOMAINS_ARGS}

echo "[SUCESSO] Certificados emitidos com sucesso!"
echo "Caminho dos certificados:"
echo "  Certificado: /etc/letsencrypt/live/${DOMAIN}/fullchain.pem"
echo "  Chave privada: /etc/letsencrypt/live/${DOMAIN}/privkey.pem"

# Copia para a pasta do Nginx se ela existir
TARGET_SSL_DIR="./nginx/ssl"
if [[ -d "${TARGET_SSL_DIR}" ]]; then
  echo "Copiando certificados para ${TARGET_SSL_DIR}..."
  cp "/etc/letsencrypt/live/${DOMAIN}/fullchain.pem" "${TARGET_SSL_DIR}/cert.pem"
  cp "/etc/letsencrypt/live/${DOMAIN}/privkey.pem" "${TARGET_SSL_DIR}/key.pem"
  chmod 600 "${TARGET_SSL_DIR}/key.pem"
  echo "[OK] Certificados atualizados na pasta nginx/ssl!"
fi

# Configuração de cron job para auto-renovação automática diária
echo "Configurando crontab para renovação automática duas vezes ao dia..."
(crontab -l 2>/dev/null || true; echo "0 0,12 * * * certbot renew --quiet --post-hook 'docker compose restart nginx || systemctl reload nginx'") | sort -u | crontab -
echo "[OK] Renovação automática via cron configurada com sucesso!"

