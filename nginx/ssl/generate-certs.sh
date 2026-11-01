#!/bin/sh
set -e

mkdir -p "$(dirname "$0")"
SSL_DIR="$(cd "$(dirname "$0")" && pwd)"

echo "Generating self-signed SSL certificates in $SSL_DIR..."

openssl req -x509 -nodes -days 3650 -newkey rsa:2048 \
  -keyout "$SSL_DIR/key.pem" \
  -out "$SSL_DIR/cert.pem" \
  -subj "/C=BR/ST=SP/L=Sao Paulo/O=ERP Agendamentos/OU=Dev/CN=localhost" \
  -addext "subjectAltName=DNS:localhost,DNS:*.localhost,IP:127.0.0.1"

echo "Certificates successfully generated:"
echo "  - $SSL_DIR/cert.pem"
echo "  - $SSL_DIR/key.pem"

