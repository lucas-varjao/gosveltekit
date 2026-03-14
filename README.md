# GoSvelteKit

Template fullstack com backend Go + Gin, frontend SvelteKit 5, autenticação baseada em sessão e referências prontas para páginas protegidas, admin e tabelas server-side.

## Criando um novo projeto

Use este repositório como GitHub Template ou clone localmente e rode o bootstrap de identidade do projeto:

```bash
git clone https://github.com/lucas-varjao/gosveltekit.git meu-novo-projeto
cd meu-novo-projeto

make init \
  APP_NAME=acme-starter \
  DISPLAY_NAME='Acme Starter' \
  GO_MODULE=github.com/acme/acme-starter \
  CONTAINER_REGISTRY=ghcr.io/acme \
  DOMAIN=acme-starter.local \
  K8S_NAMESPACE=acme-starter
```

`make init` reescreve:

- nome do app e display name
- módulo Go
- nomes e referências de imagem
- domínio público e namespace Kubernetes
- branding do frontend
- defaults de email/from
- manifesto Kubernetes e referências de documentação

Os metadados centralizados do template ficam em [`project.env`](/var/home/lvarjao/dev/pessoal/gosveltekit/project.env).

## Onboarding local

O caminho mais curto para subir um clone limpo é:

```bash
make bootstrap
make infra-up
make migrate-up
make seed-admin ADMIN_IDENTIFIER=admin ADMIN_EMAIL=admin@example.local ADMIN_PASSWORD='Starter123!'
make dev-backend
make dev-frontend
```

O `bootstrap` cria:

- [`.env`](/var/home/lvarjao/dev/pessoal/gosveltekit/.env) a partir de [`.env.example`](/var/home/lvarjao/dev/pessoal/gosveltekit/.env.example)
- [`backend/.env`](/var/home/lvarjao/dev/pessoal/gosveltekit/backend/.env) a partir de [`backend/.env.example`](/var/home/lvarjao/dev/pessoal/gosveltekit/backend/.env.example)
- [`frontend/.env`](/var/home/lvarjao/dev/pessoal/gosveltekit/frontend/.env) a partir de [`frontend/.env.example`](/var/home/lvarjao/dev/pessoal/gosveltekit/frontend/.env.example)

`infra-up` sobe:

- PostgreSQL
- Mailpit

Mailpit fica disponível por padrão em `http://localhost:8025`.

`make infra-up` e `make infra-down` usam `docker compose`, `podman compose` ou os wrappers `docker-compose` / `podman-compose`, conforme o runtime configurado em `CONTAINER_CLI`.

## Fluxos principais

Comandos relevantes do template:

- `make help`
- `make bootstrap`
- `make infra-up`
- `make infra-down`
- `make migrate-up`
- `make migrate-down`
- `make migrate-create name=create_widgets`
- `make seed-admin ADMIN_IDENTIFIER=admin ADMIN_EMAIL=admin@example.local ADMIN_PASSWORD='Starter123!'`
- `make dev-backend`
- `make dev-frontend`
- `make test`
- `make build`
- `make images`

## Migrations e seed

O runtime não executa `AutoMigrate` e não cria usuário admin implicitamente.

Migrações versionadas ficam em [`backend/db/migrations`](/var/home/lvarjao/dev/pessoal/gosveltekit/backend/db/migrations) e são aplicadas com Goose via `make migrate-up`.

Documentação detalhada: [`docs/goose-migrations.md`](/var/home/lvarjao/dev/pessoal/gosveltekit/docs/goose-migrations.md).

O usuário administrador deve ser criado explicitamente:

```bash
make seed-admin \
  ADMIN_IDENTIFIER=admin \
  ADMIN_EMAIL=admin@example.local \
  ADMIN_PASSWORD='Starter123!' \
  ADMIN_DISPLAY_NAME='Platform Admin'
```

## Build e imagens

A versão do projeto é centralizada em [`VERSION`](/var/home/lvarjao/dev/pessoal/gosveltekit/VERSION).

Build local:

```bash
make build
```

Build de imagens versionadas:

```bash
./scripts/build-images.sh
```

O script usa [`project.env`](/var/home/lvarjao/dev/pessoal/gosveltekit/project.env) para resolver nomes e refs das imagens, lê a versão de [`VERSION`](/var/home/lvarjao/dev/pessoal/gosveltekit/VERSION) e aceita:

- `CONTAINER_CLI=podman|docker`
- `VITE_API_URL=http://localhost:8080`
- `PUSH_IMAGES=true|false`

O build gera três imagens versionadas:

- backend
- frontend
- migrator

Para publicar as tags remotas configuradas em `*_IMAGE_REF` no Docker Hub ou outro registry compatível:

```bash
PUSH_IMAGES=true ./scripts/build-images.sh
```

Ou, pelo `Makefile`:

```bash
make images PUSH_IMAGES=true
```

## Kubernetes

Os manifestos base atuais são:

- [`k8s/gosveltekit-base.yaml`](/var/home/lvarjao/dev/pessoal/gosveltekit/k8s/gosveltekit-base.yaml): namespace, `ConfigMap` e `Secret`
- [`k8s/gosveltekit-migrate.job.yaml`](/var/home/lvarjao/dev/pessoal/gosveltekit/k8s/gosveltekit-migrate.job.yaml): `Job` dedicado para `goose up`
- [`k8s/gosveltekit.yaml`](/var/home/lvarjao/dev/pessoal/gosveltekit/k8s/gosveltekit.yaml): `Deployment`, `Service` e `Ingress`

Depois de `make init`, eles são renomeados para `k8s/<app-name>-base.yaml`, `k8s/<app-name>-migrate.job.yaml` e `k8s/<app-name>.yaml`.

Antes de aplicar no cluster, ajuste:

- imagens do backend e frontend
- imagem do migrator
- `DATABASE_DSN`
- credenciais SMTP
- host do ingress

Sequência recomendada de deploy:

```bash
make k8s-deploy
```

O target `make k8s-migrate-job` usa por padrão os valores de `project.env`, mas aceita override via variáveis como `KUBECTL`, `K8S_NAMESPACE`, `K8S_MIGRATE_TIMEOUT`, `K8S_BASE_MANIFEST` e `K8S_MIGRATE_JOB_MANIFEST`.

Se você quiser executar rollout completo, `make k8s-deploy` chama `make k8s-migrate-job`, aplica `k8s/<app>.yaml` e aguarda o rollout de backend e frontend. Ele também aceita override de `K8S_APP_MANIFEST`, `K8S_BACKEND_DEPLOYMENT_NAME`, `K8S_FRONTEND_DEPLOYMENT_NAME` e `K8S_ROLLOUT_TIMEOUT`.

Para operar frontend e backend sob o mesmo host no ingress, gere a imagem do frontend com `VITE_API_URL=''`.

## O que é core e o que é opcional

Partes core do starter:

- autenticação baseada em sessão
- páginas de login, registro e conta
- backend Gin + GORM + PostgreSQL
- frontend SvelteKit 5 + Tailwind + shadcn-svelte
- versionamento central via [`VERSION`](/var/home/lvarjao/dev/pessoal/gosveltekit/VERSION)
- build de imagens e manifesto Kubernetes

Partes opcionais/removíveis:

- hub administrativo e exemplos em `/admin`
- playground de paginação em `/examples/pagination`
- dados mockados para exemplos de tabela

## Estrutura

```text
.
├── backend/
│   ├── cmd/
│   │   ├── migrate/
│   │   └── seed-admin/
│   ├── db/migrations/
│   ├── internal/
│   └── main.go
├── frontend/
├── k8s/
├── scripts/
├── project.env
└── VERSION
```

## Qualidade do template

O repositório inclui dois fluxos de CI:

- `quality`: executa `make test` e `make build`
- `template-smoke`: cria um diretório temporário, roda `make init`, sobe infra, aplica migração, faz seed de admin, valida build/login e testa o build de imagens
