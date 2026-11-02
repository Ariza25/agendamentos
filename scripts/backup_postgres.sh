#!/usr/bin/env bash
# =====================================================================
# Script de Backup Automatizado do PostgreSQL (ERP Agendamentos)
# =====================================================================
# Realiza pg_dump do banco de dados, comprime com gzip, remove backups
# com mais de 7 dias de retenção e exibe status da operação.
#
# Uso:
#   export DATABASE_URL="postgres://user:pass@host:5432/dbname?sslmode=require"
#   ./scripts/backup_postgres.sh
# =====================================================================

set -eo pipefail

BACKUP_DIR="${BACKUP_DIR:-./backups/postgres}"
RETENTION_DAYS="${RETENTION_DAYS:-7}"
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")
BACKUP_FILE="${BACKUP_DIR}/erpagendamentos_backup_${TIMESTAMP}.sql.gz"

mkdir -p "${BACKUP_DIR}"

if [[ -z "${DATABASE_URL}" ]]; then
  echo "[ERRO] A variável de ambiente DATABASE_URL não foi informada."
  echo "Exemplo: export DATABASE_URL='postgresql://user:pass@host:5432/dbname?sslmode=require'"
  exit 1
fi

echo "=========================================================="
echo "Iniciando backup do banco de dados: ${TIMESTAMP}"
echo "Destino: ${BACKUP_FILE}"
echo "=========================================================="

# Executa pg_dump e comprime diretamente para gzip
pg_dump "${DATABASE_URL}" \
  --no-owner \
  --no-acl \
  --clean \
  --if-exists \
  | gzip -9 > "${BACKUP_FILE}"

FILESIZE=$(du -h "${BACKUP_FILE}" | cut -f1)
echo "[SUCESSO] Backup concluído com êxito! Tamanho: ${FILESIZE}"

# Limpeza de backups antigos conforme a política de retenção
echo "Aplicando política de retenção: removendo arquivos com mais de ${RETENTION_DAYS} dias em ${BACKUP_DIR}..."
find "${BACKUP_DIR}" -type f -name "erpagendamentos_backup_*.sql.gz" -mtime +"${RETENTION_DAYS}" -exec rm -vf {} \;

echo "[OK] Rotina de backup finalizada."

