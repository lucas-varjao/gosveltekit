# Migrações com Goose

Esta base usa migrações versionadas com [`goose`](https://github.com/pressly/goose) em vez de `AutoMigrate` no startup da aplicação.

O objetivo é simples:

- ter histórico explícito de schema
- controlar a ordem das mudanças
- separar startup da API de alteração estrutural no banco
- tornar deploy em produção previsível

## Como está organizado

As peças principais são:

- comando de migração: [`backend/cmd/migrate/main.go`](/var/home/lvarjao/dev/pessoal/gosveltekit/backend/cmd/migrate/main.go)
- diretório de migrações: [`backend/db/migrations`](/var/home/lvarjao/dev/pessoal/gosveltekit/backend/db/migrations)
- conexão SQL para o Goose: [`backend/internal/bootstrap/database.go`](/var/home/lvarjao/dev/pessoal/gosveltekit/backend/internal/bootstrap/database.go)
- targets do Makefile: [`Makefile`](/var/home/lvarjao/dev/pessoal/gosveltekit/Makefile)

O comando `backend/cmd/migrate` suporta:

- `up`: aplica migrações pendentes
- `down`: desfaz a última migração aplicada
- `create <name>`: cria um novo arquivo `.sql`

## Como o comando funciona

O fluxo do [`backend/cmd/migrate/main.go`](/var/home/lvarjao/dev/pessoal/gosveltekit/backend/cmd/migrate/main.go) é:

1. carrega a configuração com `config.LoadConfig()`
2. abre conexão `database/sql` via `pgx`
3. configura o dialect `postgres` no Goose
4. executa `goose.Up(...)`, `goose.Down(...)` ou cria um novo arquivo de migração

O DSN vem do mesmo lugar que a aplicação:

- `DATABASE_DSN`
- ou `DATABASE_URL`
- ou fallback do `backend/configs/app.yml`

Isso garante que migrator e aplicação apontem para o mesmo banco.

## Estrutura dos arquivos de migração

Cada arquivo fica em `backend/db/migrations` e usa blocos `Up` e `Down`.

Exemplo real do projeto:

[`backend/db/migrations/20260314120000_create_users_and_sessions.sql`](/var/home/lvarjao/dev/pessoal/gosveltekit/backend/db/migrations/20260314120000_create_users_and_sessions.sql)

Estrutura:

```sql
-- +goose Up
-- +goose StatementBegin
CREATE TABLE ...;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE ...;
-- +goose StatementEnd
```

Regras práticas:

- `Up` deve levar o schema para o estado novo
- `Down` deve reverter exatamente o que o `Up` criou, quando isso for seguro
- dê nomes objetivos ao arquivo: `create_users`, `add_user_status`, `create_audit_log`
- trate mudanças destrutivas com mais cuidado; rollback nem sempre é seguro em produção

## Fluxo local

### 1. Subir infra

```bash
make bootstrap
make infra-up
```

### 2. Aplicar migrações

```bash
make migrate-up
```

Isso executa:

```bash
cd backend
go run ./cmd/migrate up
```

### 3. Criar uma nova migração

```bash
make migrate-create name=add_user_status
```

Isso gera um arquivo como:

```text
backend/db/migrations/20260314153000_add_user_status.sql
```

### 4. Reverter a última migração

```bash
make migrate-down
```

Use rollback com cuidado. Em ambiente compartilhado ou produção, `down` só deve ser usado com plano claro de impacto.

## Relação com o startup da API

O binário principal do backend em [`backend/main.go`](/var/home/lvarjao/dev/pessoal/gosveltekit/backend/main.go) não faz:

- `AutoMigrate`
- criação automática de tabelas
- bootstrap implícito de admin

Ou seja: a API assume que o schema já está pronto antes de subir.

Se o banco não estiver migrado, o deploy pode até iniciar o processo, mas a aplicação vai falhar ao acessar estruturas que ainda não existem.

## Como isso funciona no deploy

Em produção, a ordem correta é:

1. publicar a nova imagem
2. executar as migrações no banco alvo
3. só depois liberar o rollout da nova versão da API

Essa ordem evita que pods novos e antigos disputem schema incompatível.

## Estado atual do projeto no Kubernetes

O manifesto atual [`k8s/gosveltekit.yaml`](/var/home/lvarjao/dev/pessoal/gosveltekit/k8s/gosveltekit.yaml) sobe:

- namespace
- configmap
- secret
- deployment/service do backend
- deployment/service do frontend
- ingress

Hoje ele **não** inclui:

- `Job` de migração
- `initContainer` para migrar schema
- `helm hook` ou mecanismo equivalente

Além disso, a imagem final do backend em [`backend/Dockerfile`](/var/home/lvarjao/dev/pessoal/gosveltekit/backend/Dockerfile) copia apenas o binário `server`.

Na prática, isso significa que o deploy atual em Kubernetes depende de um passo externo para rodar migração.

## Estratégias recomendadas para Kubernetes

### Opção 1: migrar antes do `kubectl apply`

É a estratégia mais simples e segura para este repositório no estado atual.

Fluxo:

1. CI builda e publica a imagem nova
2. CI executa `go run ./cmd/migrate up` contra o banco de produção
3. CI aplica o manifesto do backend/frontend

Vantagens:

- simples
- explícito
- evita acoplamento da migração ao ciclo de vida do pod

Cuidados:

- o job de CI precisa acesso seguro ao banco
- mudanças incompatíveis entre versões precisam ser backward compatible durante rollout

## Opção 2: Job de migração no cluster

É a opção mais comum quando você quer que o cluster execute a migração.

Neste template, essa opção agora está implementada com:

- uma imagem `migrator` dedicada, buildada a partir de [`backend/Dockerfile`](/var/home/lvarjao/dev/pessoal/gosveltekit/backend/Dockerfile)
- um manifesto de `Job` em [`k8s/gosveltekit-migrate.job.yaml`](/var/home/lvarjao/dev/pessoal/gosveltekit/k8s/gosveltekit-migrate.job.yaml)
- um manifesto base em [`k8s/gosveltekit-base.yaml`](/var/home/lvarjao/dev/pessoal/gosveltekit/k8s/gosveltekit-base.yaml) para `Namespace`, `ConfigMap` e `Secret`

A imagem `migrator` contém:

- o binário `./migrate`
- o diretório `db/migrations`
- a mesma pasta `configs` usada pela aplicação

O comando default do container já executa `./migrate up`, e o `Job` também declara isso explicitamente.

Fluxo operacional:

1. aplicar [`k8s/gosveltekit-base.yaml`](/var/home/lvarjao/dev/pessoal/gosveltekit/k8s/gosveltekit-base.yaml)
2. recriar o `Job` de migração
3. aguardar sucesso do job
4. aplicar [`k8s/gosveltekit.yaml`](/var/home/lvarjao/dev/pessoal/gosveltekit/k8s/gosveltekit.yaml)

Exemplo prático:

```bash
make k8s-deploy
```

O target resolve por padrão:

- `K8S_NAMESPACE` a partir de `project.env`
- `K8S_BASE_MANIFEST` como `k8s/<app-slug>-base.yaml`
- `K8S_MIGRATE_JOB_MANIFEST` como `k8s/<app-slug>-migrate.job.yaml`
- `K8S_MIGRATE_JOB_NAME` como `<app-slug>-migrate`

O target `k8s-deploy` também resolve por padrão:

- `K8S_APP_MANIFEST` como `k8s/<app-slug>.yaml`
- `K8S_BACKEND_DEPLOYMENT_NAME` como `<app-slug>-backend`
- `K8S_FRONTEND_DEPLOYMENT_NAME` como `<app-slug>-frontend`

Se necessário, você pode sobrescrever esses valores na invocação do `make`.

Vantagens:

- a migração roda dentro do cluster
- usa a mesma rede e os mesmos secrets do app
- separa claramente recursos base, migração e rollout da aplicação

Cuidados:

- `Job` é recurso imutável em pontos importantes; por isso o fluxo usa `delete` + `create`
- não é ideal acoplar migração ao startup de cada pod
- o rollout do app deve depender do sucesso do job

## Opção 3: initContainer no backend

Tecnicamente possível, mas normalmente não é a melhor opção.

Problema principal:

- com múltiplas réplicas, vários pods podem tentar rodar migração ao mesmo tempo

Mesmo que o Goose tenha controle de versão, esse modelo adiciona contenção e aumenta a chance de comportamento operacional confuso.

Para este template, a recomendação é:

- preferir pipeline pré-deploy
- ou um `Job` dedicado

Não usar `initContainer` como padrão.

## Recomendação para este projeto

Para o estado atual do template, a abordagem recomendada é:

### Desenvolvimento/local

```bash
make infra-up
make migrate-up
make dev-backend
make dev-frontend
```

### Produção/Kubernetes

1. buildar e publicar as imagens novas, incluindo `migrator`
2. aplicar `k8s/gosveltekit-base.yaml`
3. rodar `k8s/gosveltekit-migrate.job.yaml`
4. só depois aplicar `k8s/gosveltekit.yaml`

Esse já é o fluxo suportado pelo template.

## Exemplo de pipeline de deploy

Exemplo conceitual:

```bash
./scripts/build-images.sh

# publicar imagens no registry

make k8s-deploy
```

## Boas práticas de schema change

- prefira migrações pequenas
- evite mudanças destrutivas no mesmo deploy em que o código novo entra
- em alterações grandes, faça rollout em duas etapas
- garanta compatibilidade entre código novo e schema antigo durante a janela de rollout
- só use `down` em produção quando o plano de rollback estiver testado

## Resumo

Neste template, Goose é o mecanismo oficial de evolução de banco.

O app não migra schema sozinho; a migração continua sendo uma responsabilidade operacional explícita.

No Kubernetes, o template agora suporta isso com uma imagem `migrator` e um `Job` dedicado executado antes do rollout da aplicação.
