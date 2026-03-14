## Project Overview

GoSvelteKit is a template/base project with a Golang backend and SvelteKit (Svelte 5) frontend, designed to kickstart new projects with authentication and basic functionality already implemented.

## Backend Specifications

- **Framework**: Golang with Gin
- **Authentication**: Session-based (pluggable adapters inspired by Lucia Auth)
- **Database**: PostgreSQL em runtime (via GORM)
- **ORM**: GORM for database operations
- **API Docs**: Swaggo
- **Logging**: slog (log/slog)
- **Config**: Viper com padrão env-first e fallback para `backend/configs/app.yml`
- **Versioning**: versão centralizada no arquivo raiz `VERSION`

### Config Loading Pattern

- Sempre priorizar variáveis de ambiente para todas as chaves de configuração.
- Quando a variável não existir, usar o valor do arquivo `backend/configs/app.yml`.
- Convenção de nomes: chaves aninhadas viram env em uppercase com `_`:
    - `server.port` -> `SERVER_PORT`
    - `auth.session_ttl` -> `AUTH_SESSION_TTL`
    - `email.smtp_host` -> `EMAIL_SMTP_HOST`
- Banco: `DATABASE_DSN` é o nome preferencial; `DATABASE_URL` é alias compatível.

### Authentication Architecture

The auth system uses a pluggable adapter pattern:

```
internal/auth/
├── interfaces.go      # UserAdapter, SessionAdapter interfaces
├── auth_manager.go    # Central AuthManager
└── adapter/gorm/      # GORM implementation
```

- **Sessions** are stored in the database (not JWTs)
- Login returns `session_id` (not `access_token`/`refresh_token`)
- Auth via `Authorization: Bearer {session_id}` header or `session_id` cookie

## Frontend Specifications

- **Framework**: SvelteKit (Svelte 5)
- **Styling**: TailwindCSS (utility-first)
- **Runtime**: Bun
- **UI Components**: shadcn-svelte
- **Icons**: @lucide/svelte (Svelte 5)
- **Mode**: Dark mode only
- Always use Svelte 5 with the new runes API (`$state`, `$derived`, `$props`)
- Data tables should use `shadcn-svelte` + `TanStack Table`
- Prefer server-side tables: pagination, filter, and sorting must come from the Go backend
- The reference implementation for this stack lives in `frontend/src/routes/(protected)/admin/+page.svelte`
- A versão exibida no footer deve vir da versão central do projeto

### Auth Store

```typescript
// User fields are snake_case
interface User {
    id: string;
    identifier: string; // username
    email: string;
    display_name: string;
    role: string;
    active: boolean;
}
```

## Design Preferences

- Professional and minimalist design
- Basic animations for enhanced UX
- `slate-950` as primary background color

## Code Style & Conventions

- Backend Go code follows standard Go idioms and project layout best practices
- Frontend TypeScript uses 4-space indentation and single quotes
- CSS is primarily written as Tailwind utility classes
- Component files use .svelte extension and follow SvelteKit conventions

## Development Workflow

- Monorepo with `backend/` and `frontend/` directories
- Dev requires running both servers: `cd backend && go run main.go` and `cd frontend && bun run dev`
- Existe um `Makefile` na raiz para tarefas comuns (`make help`, `make test`, `make images`, etc.)
- Follow Conventional Commits for commit messages
- Ao alterar a versão do projeto, atualizar apenas o arquivo raiz `VERSION`

## Deployment

- Docker containers for deployment
- Runtime final previsto: Kubernetes
- Environment variables for configuration
- Imagens devem ser versionadas a partir de `VERSION`
- O script padrão para build das imagens é `./scripts/build-images.sh`
- O script usa `podman` por padrão e aceita `CONTAINER_CLI=docker`
- Existe um manifesto único em `k8s/gosveltekit.yaml`
- As imagens de runtime usam `TZ=America/Sao_Paulo`
