# GoSvelteKit

Um template fullstack pronto para uso com autenticação JWT, combinando backend Golang com SQLite e frontend SvelteKit.

## 📋 Visão Geral

GoSvelteKit é um projeto base projetado para acelerar o desenvolvimento de aplicações web fullstack. Este template vem pré-configurado com autenticação JWT, banco de dados SQLite e páginas de login/registro, permitindo que você pule a configuração inicial repetitiva e foque nas funcionalidades específicas do seu projeto.

## 🚀 Recursos

### Backend (Golang)

-   Sistema de autenticação JWT completo
-   Banco de dados SQLite com migrations
-   Estrutura modular e escalável
-   Middleware de autenticação
-   API RESTful

### Frontend (SvelteKit)

-   Páginas de autenticação prontas (login, registro, recuperação de senha)
-   Gerenciamento de estado para autenticação
-   Layout responsivo básico
-   Interceptação de requisições para inclusão de tokens

## 🛠️ Pré-requisitos

-   Go 1.18+
-   Node.js 16+
-   npm ou pnpm
-   Docker e Docker Compose (opcional)

## 🔧 Instalação e Uso

### Clonando o template

```bash
git clone https://github.com/seu-usuario/gosveltekit.git meu-novo-projeto
cd meu-novo-projeto
```

### Usando Docker Compose (recomendado)

```bash
docker-compose up
```

### Execução manual

#### Backend

```bash
cd backend
go mod download
go run cmd/api/main.go
```

#### Frontend

```bash
cd frontend
npm install
npm run dev
```

## 📁 Estrutura do Projeto

```
gosveltekit/
├── backend/               # Aplicação Golang
│   ├── cmd/api/           # Ponto de entrada
│   ├── internal/          # Código principal
│   │   ├── auth/          # Autenticação JWT
│   │   ├── config/        # Configurações
│   │   ├── handlers/      # Controladores HTTP
│   │   ├── middleware/    # Middlewares
│   │   ├── models/        # Definições de dados
│   │   └── repository/    # Acesso ao banco
│   └── migrations/        # Migrações SQL
│
└── frontend/              # Aplicação SvelteKit
    ├── src/
    │   ├── lib/           # Utilitários e componentes
    │   └── routes/        # Páginas da aplicação
    │       ├── auth/      # Rotas de autenticação
    │       └── dashboard/ # Área autenticada
```

## ⚙️ Configuração

### Variáveis de Ambiente

Copie o arquivo `.env.example` para `.env` e ajuste as variáveis conforme necessário:

```bash
cp .env.example .env
```

### Banco de Dados

As migrações são executadas automaticamente na inicialização. Para criar novas migrações:

```bash
cd backend
go run cmd/migrate/main.go create "nome_da_migracao"
```

## 🔄 Começando um Novo Projeto

1. Clone este repositório com um novo nome
2. Personalize o `.env` e as configurações
3. Modifique os modelos no backend conforme necessário
4. Adapte as páginas do frontend para seu caso de uso
5. Inicie o desenvolvimento das funcionalidades específicas do seu projeto

## 📄 Licença

Este projeto está licenciado sob a MIT License - veja o arquivo [LICENSE](LICENSE) para detalhes.

## 🤝 Contribuição

Contribuições são bem-vindas! Por favor, sinta-se à vontade para enviar um pull request.

---

Desenvolvido com ❤️ para agilizar seu fluxo de trabalho de desenvolvimento.
