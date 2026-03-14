# GoSvelteKit

Um template fullstack pronto para uso com **autenticação baseada em sessões**, combinando backend Golang com PostgreSQL e frontend SvelteKit.

## 📋 Visão Geral

GoSvelteKit é um projeto base projetado para acelerar o desenvolvimento de aplicações web fullstack. Este template vem pré-configurado com autenticação plugável (inspirada no Lucia Auth), banco de dados PostgreSQL em runtime e páginas de login/registro, permitindo que você pule a configuração inicial repetitiva e foque nas funcionalidades específicas do seu projeto.

## 🚀 Recursos

### Backend (Golang)

- **Autenticação plugável** com adapters (estilo Lucia Auth)
- Sessões armazenadas no banco de dados
- Banco de dados PostgreSQL com GORM (runtime)
- Estrutura modular e escalável
- Middleware de autenticação
- API RESTful com Gin

### Frontend (SvelteKit)

- Páginas de autenticação prontas (login, registro, recuperação de senha)
- Gerenciamento de estado com Svelte 5 runes (`$state`, `$derived`)
- Layout responsivo com TailwindCSS
- Componentes UI seguindo o padrão **shadcn-svelte**
- Exemplo de Data Table server-side em `/admin` com **TanStack Table**
- Ícones SVG com **@lucide/svelte** (Svelte 5)
- Sessão baseada em cookie HttpOnly no navegador

## 🛠️ Pré-requisitos

- Go 1.26+
- Bun (ou Node.js 18+)
- Podman ou Docker (opcional)

## 🔧 Instalação e Uso

### Clonando o template

```bash
git clone https://github.com/lucas-varjao/gosveltekit.git meu-novo-projeto
cd meu-novo-projeto
```

### Atalhos com Makefile

```bash
make help
make test
make images
```

Alvos úteis:

- `make dev-backend`
- `make dev-frontend`
- `make build`
- `make test`
- `make images`

### Build versionado de imagens

```bash
./scripts/build-images.sh
```

Por padrão o script usa `podman`, lê a versão de `VERSION` e gera as tags:

- `gosveltekit-backend:<versao>` e `gosveltekit-backend:latest`
- `gosveltekit-frontend:<versao>` e `gosveltekit-frontend:latest`

As imagens de runtime usam o fuso horário `America/Sao_Paulo`.

Para usar Docker em vez de Podman:

```bash
CONTAINER_CLI=docker ./scripts/build-images.sh
```

Para apontar o frontend para outra API no build da imagem:

```bash
VITE_API_URL='https://api.seu-dominio.com' ./scripts/build-images.sh
```

### Execução manual

#### Backend

```bash
cd backend
go mod download
go run main.go
```

Ao iniciar localmente via `backend/`, o servidor resolve a versão a partir do arquivo raiz `../VERSION` e exibe esse valor no log de startup.

#### Frontend

```bash
cd frontend
bun install
bun run dev
```

Em desenvolvimento e build de produção, o frontend lê a mesma versão central de `../VERSION` e a exibe no rodapé.

Você também pode usar os alvos do `Makefile`:

```bash
make dev-backend
make dev-frontend
```

### Execução das imagens

#### Backend

```bash
podman run --rm -p 8080:8080 \
  -e DATABASE_DSN='postgresql://gosvelte:gosvelte@host.containers.internal:5432/gosveltekit?sslmode=disable' \
  gosveltekit-backend:0.1.0
```

Com Docker, ajuste o host do banco conforme o seu ambiente.

#### Frontend

```bash
podman run --rm -p 3000:80 gosveltekit-frontend:0.1.0
```

## ☸️ Kubernetes

O projeto inclui um manifesto único em `k8s/gosveltekit.yaml`.

Ele contém:

- `Namespace`
- `ConfigMap`
- `Secret`
- `Deployment` e `Service` do backend
- `Deployment` e `Service` do frontend
- `Ingress`

Antes de aplicar no cluster, ajuste:

- as imagens `ghcr.io/your-org/gosveltekit-backend:0.1.0` e `ghcr.io/your-org/gosveltekit-frontend:0.1.0`
- o `DATABASE_DSN` e credenciais SMTP no `Secret`
- o host do ingress (`gosveltekit.local`)

Aplicação:

```bash
kubectl apply -f k8s/gosveltekit.yaml
```

Para o modelo de ingress do manifesto funcionar sem origem cruzada, gere a imagem do frontend com `VITE_API_URL=''` ou com a URL pública real da API.

## 📁 Estrutura do Projeto

```bash
gosveltekit/
├── backend/
│   ├── main.go               # Ponto de entrada
│   └── internal/
│       ├── auth/             # Sistema de autenticação
│       │   ├── interfaces.go # UserAdapter, SessionAdapter
│       │   ├── auth_manager.go
│       │   └── adapter/gorm/ # Implementação GORM
│       ├── config/
│       ├── handlers/
│       ├── middleware/
│       ├── models/
│       ├── repository/
│       ├── router/
│       └── service/
│
└── frontend/
    └── src/
        ├── lib/
        │   ├── api/          # Cliente HTTP, auth e listagens administrativas
        │   └── stores/       # Estado (auth store)
        └── routes/
            ├── login/
            ├── register/
            └── (protected)/  # Rotas autenticadas
```

## 📊 Data Tables

O projeto agora inclui uma referência oficial de tabela administrativa em `frontend/src/routes/(protected)/admin/+page.svelte`.

Direção da stack:

- frontend com **Svelte 5** + **shadcn-svelte**
- comportamento de tabela com **TanStack Table**
- paginação, filtro e sorting controlados pelo backend Go
- contrato HTTP explícito para listagens server-side

### Endpoint de referência

```http
GET /api/admin/users?page=1&page_size=10&search=admin&sort=created_at&order=desc
```

### Resposta de referência

```json
{
    "items": [
        {
            "id": "1",
            "identifier": "admin",
            "email": "admin@example.com",
            "display_name": "Administrator",
            "role": "admin",
            "active": true,
            "last_login": "2026-03-14T10:00:00Z",
            "created_at": "2026-01-10T08:00:00Z"
        }
    ],
    "page": 1,
    "page_size": 10,
    "total_items": 1,
    "total_pages": 1,
    "sort": {
        "field": "created_at",
        "direction": "desc"
    },
    "search": "admin"
}
```

### Convenção recomendada

Para novas listagens administrativas:

- crie endpoints paginados no backend Go em vez de retornar arrays completos
- mantenha filtros e ordenação como query params simples
- use `snake_case` no contrato JSON
- trate `TanStack Table` como camada de comportamento, e `shadcn-svelte` como base visual
- use a tela `/admin` como referência para paginação, busca e ordenação server-side

## 🔐 Autenticação

O sistema usa **autenticação baseada em sessões** com adapters plugáveis:

```go
// Interfaces que você pode implementar para qualquer banco
type UserAdapter interface {
    FindUserByIdentifier(identifier string) (*UserData, error)
    ValidateCredentials(identifier, password string) (*UserData, error)
    // ...
}

type SessionAdapter interface {
    CreateSession(userID string, expiresAt time.Time, metadata SessionMetadata) (*Session, error)
    GetSession(sessionID string) (*Session, error)
    // ...
}
```

### Resposta de Login

```json
{
    "session_id": "abc123...",
    "expires_at": "2024-02-11T12:00:00Z",
    "user": {
        "id": "1",
        "identifier": "admin",
        "email": "admin@example.com",
        "display_name": "Administrator",
        "role": "admin"
    }
}
```

### Canais de autenticação suportados

- Web: cookie `session_id` (HttpOnly)
- API clients/mobile/CLI: `Authorization: Bearer {session_id}` ou `X-Session-ID`

### Exemplos cURL (CLI)

```bash
# 1) Login e captura do session_id
SESSION_ID=$(curl -s -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin"}' | jq -r '.session_id')

echo "$SESSION_ID"
```

```bash
# 2) Acesso via Authorization: Bearer
curl -s http://localhost:8080/api/me \
  -H "Authorization: Bearer ${SESSION_ID}"
```

```bash
# 3) Acesso via X-Session-ID
curl -s http://localhost:8080/api/me \
  -H "X-Session-ID: ${SESSION_ID}"
```

```bash
# 4) Fluxo por cookie (estilo navegador)
# Salva cookies após login
curl -s -c cookies.txt -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin"}'

# Reusa cookie para rota protegida
curl -s -b cookies.txt http://localhost:8080/api/me
```

## ⚙️ Configuração

Copie o arquivo `.env.example` para `.env` e ajuste as variáveis conforme necessário:

```bash
cp .env.example .env
```

O backend usa **Viper com padrão env-first**:

- Se a variável de ambiente existir, ela vence.
- Se não existir, faz fallback para `backend/configs/app.yml`.

Isso vale para todas as seções (`server`, `database`, `auth`, `email`) usando o padrão `SECAO_CHAVE`.

Exemplos:

```bash
export SERVER_PORT='8080'
export AUTH_SESSION_TTL='720h'
export EMAIL_SMTP_HOST='sandbox.smtp.mailtrap.io'
export DATABASE_DSN='postgresql://postgres:postgres@localhost:5432/gosveltekit?sslmode=disable'
```

Compatibilidade de banco: `DATABASE_URL` também é aceito como alias de `DATABASE_DSN`.

Observação sobre testes: a suíte automatizada do backend usa SQLite em memória para manter execução rápida.

## 🏷️ Versionamento

O arquivo raiz `VERSION` é a fonte canônica da versão do projeto.

```bash
cat VERSION
0.1.0
```

Regras:

- O formato é `MAJOR.MINOR.PATCH`, sem prefixo `v`
- Backend, frontend e imagens devem consumir esse mesmo valor
- A interface e os logs podem exibir `v` ao renderizar a versão

Fluxo:

- Backend: usa `ldflags`, depois `APP_VERSION`, depois `../VERSION`, e por fim fallback `dev`
- Frontend: usa `APP_VERSION` no build de imagem e `../VERSION` em dev/build local
- Imagens: recebem `APP_VERSION` no build e publicam label OCI `org.opencontainers.image.version`
- Tags geradas pelo script: `gosveltekit-backend:<versao>`, `gosveltekit-backend:latest`, `gosveltekit-frontend:<versao>` e `gosveltekit-frontend:latest`

## 🔄 Começando um Novo Projeto

1. Clone este repositório com um novo nome
2. Personalize o `.env` e as configurações
3. Modifique os modelos no backend conforme necessário
4. Adapte as páginas do frontend para seu caso de uso
5. Para integrar com outro banco de usuários, implemente `UserAdapter`

## 📄 Licença

Este projeto está licenciado sob a MIT License - veja o arquivo [LICENSE](LICENSE) para detalhes.

## 🤝 Contribuição

Contribuições são bem-vindas! Por favor, sinta-se à vontade para enviar um pull request.

---

Desenvolvido com ❤️ para agilizar seu fluxo de trabalho de desenvolvimento.
