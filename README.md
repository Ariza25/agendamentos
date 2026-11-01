# ERP Agendamentos

Sistema de agendamentos online: catálogo de serviços, profissionais, clientes e agenda de
atendimentos em uma única plataforma, com página pública de agendamento para cada
estabelecimento (multi-tenant).

## Visão geral

- **Backend**: API REST em Go, organizada em camadas (`domain`, `application`,
  `interfaces/http/rest`, `infrastructure`, `bootstrap`), multi-tenant com isolamento por
  `tenant_id`.
- **Admin** (`admin/`): painel de gestão (Vite + React + Tailwind) para configuração de
  serviços, profissionais, clientes, agendamentos, configuração e tenants (somente master).
- **Booking** (`booking/`): página pública de agendamento (Vite + React + Tailwind) exposta
  por slug de estabelecimento.
- **Landing** (`landing/`): página de apresentação do produto.

## Stack

| Componente | Tecnologia |
| ---------- | ---------- |
| Backend | Go 1.25, PostgreSQL 16 |
| Admin / Booking / Landing | Vite + React 18 + Tailwind CSS 4 |

## Estrutura

```
erpagendamentos/
├── backend/          # API REST em Go
│   ├── cmd/api/      # entrypoint
│   ├── internal/
│   │   ├── domain/           # entidades e contratos
│   │   ├── application/      # regras e orquestração (use cases)
│   │   ├── interfaces/http/rest/  # rotas, handlers, auth
│   │   ├── infrastructure/   # Postgres, memória, segurança
│   │   └── bootstrap/        # wiring e composition root
│   └── sql/          # schema unificado (schema.sql)
├── admin/            # painel de gestão
├── booking/          # página pública de agendamento
├── landing/          # landing page do produto
├── docker-compose.yml
└── .github/workflows/ci.yml
```

## Pré-requisitos

- Go 1.25+
- Node.js 20+
- Docker (com Docker Compose) para Postgres e produção
- PostgreSQL 16 (alternativa ao Docker)

## Rodando em desenvolvimento (sem banco)

O backend usa repositórios em memória quando `ERPAGENDAMENTOS_ENV` não é `production`.

```bash
# Backend (porta 8080)
cd backend
go run ./cmd/api

# Admin (porta 5173)
cd admin
npm install && npm run dev

# Booking (porta 5175)
cd booking
npm install && npm run dev

# Landing (porta 5174)
cd landing
npm install && npm run dev
```

Em desenvolvimento, `ADMIN_EMAIL` e `ADMIN_PASSWORD` opcionais fazem o seed de um usuário
admin. O tenant `default` é criado automaticamente.

## Rodando com Docker (produção local)

```bash
docker compose up -d --build
```

Isso sobe o Postgres e o backend (porta 8080). As migrations são aplicadas automaticamente
na inicialização. Para os frontends, rode os passos de desenvolvimento acima ou sirva os
builds estáticos separadamente.

## Variáveis de ambiente (backend)

| Variável | Descrição | Padrão |
| -------- | --------- | ------ |
| `PORT` | Porta HTTP | `8080` |
| `ERPAGENDAMENTOS_ENV` | `production` ativa Postgres e validações de produção | vazio (dev) |
| `DATABASE_URL` | URL do Postgres (obrigatória em produção) | — |
| `JWT_SECRET` | Segredo JWT (obrigatório em produção, não pode ser o default) | — |
| `MFA_ENCRYPTION_KEY` | Chave de criptografia do MFA (obrigatória em produção) | — |
| `FEATURE_PRIVILEGED_MFA_ENFORCEMENT` | Exige MFA para mutações sensíveis | `true` |
| `ADMIN_EMAIL` / `ADMIN_PASSWORD` | Credenciais do seed de admin | — |
| `ALLOWED_ORIGINS` | Origens permitidas no CORS (além dos localhosts de dev) | — |
| `TZ` | Fuso horário | `America/Sao_Paulo` |

No `docker-compose.yml`, `JWT_SECRET`, `MFA_ENCRYPTION_KEY` e `ADMIN_PASSWORD` têm default
apenas para ambiente local — **troque antes de expor o serviço**. Para desenvolvimento
rápido com MFA desabilitado, defina `FEATURE_PRIVILEGED_MFA_ENFORCEMENT=false`.

## Endpoints principais

**Público** (`/api/v1/public/{slug}`)
- `GET /info`, `GET /services`, `GET /professionals`, `GET /availability`
- `POST /appointments`

**Autenticação** (`/api/v1/auth`)
- `POST /login`, `POST /logout`, `GET /me`
- MFA: `GET/POST/DELETE /mfa`, `POST /mfa/step-up`
- Usuários admin: `GET/POST /api/v1/admin/users`

**Admin** (`/api/v1/admin`, protegido por auth + permissão)
- `tenants`, `services`, `professionals`, `customers`, `appointments`, `configuration`

**Health**: `GET /health`, `GET /ready`

## Banco de Dados & Schema

- `backend/sql/schema.sql` — schema canônico único com todas as tabelas, índices e dados iniciais aplicados de forma idempotente na inicialização.

## CI

O workflow `.github/workflows/ci.yml` executa `go vet`, `go build` e `go test` no backend e
`npm run build` nos apps admin, booking e landing.

## Roles e permissões

- `owner` / `admin` — acesso administrativo ao tenant.
- `receptionist` / `professional` — acesso operacional (staff).
- Master admin (`tenantId` `default`) — acesso à área de tenants (multi-tenant).
- Mutações sensíveis exigem MFA quando o enforcement está habilitado.

## Status

Funcionalidades disponíveis: agenda e agendamentos, agendamento online público, serviços,
profissionais, clientes, configuração do estabelecimento e gestão de tenants (master).
