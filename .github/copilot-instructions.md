# Copilot Instructions for GoSvelteKit

## Project Overview

GoSvelteKit is a template/base project with a Golang backend and SvelteKit frontend, designed to kickstart new projects with authentication and basic functionality already implemented.

## Backend Specifications

-   We use Golang with Gin for our web framework
-   Authentication is handled via JWT tokens
-   SQLite is our database of choice for simplicity and portability
-   We utilize GORM as our ORM for database operations
-   API documentation is generated using Swaggo
-   Structured logging is implemented with Zap
-   Configuration management is handled by Viper

## Frontend Specifications

-   SvelteKit is our frontend framework
-   We use TailwindCSS for styling with a utility-first approach
-   Bun is our JavaScript runtime and package manager
-   UI components are built with shadcn-svelte

## Code Style & Conventions

-   Backend Go code follows standard Go idioms and the project structure adheres to Go project layout best practices
-   Frontend JavaScript/TypeScript uses 4-space indentation and single quotes
-   CSS is primarily written as Tailwind utility classes
-   Component files use .svelte extension and follow SvelteKit conventions

## Development Workflow

-   The project uses a monorepo approach with backend and frontend in separate directories
-   Development requires running both backend and frontend servers concurrently
-   We follow Conventional Commits for commit messages

## Deployment

-   The project is designed to be easily deployable as Docker containers
-   Environment variables are used for configuration across different environments
