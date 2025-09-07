# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Communication Language

**IMPORTANT: Always respond in Japanese (日本語) when working in this repository.** All explanations, error messages, and communication should be in Japanese to match the project's primary language and documentation.

## Project Overview

Mosscat is an OSS (Open Source Software) management system for internal company use, built as a full-stack application with Go backend and Vue.js frontend. It manages OSS components, versions, projects, and usage tracking with audit logging capabilities.

## Architecture

### Backend (Go)
- **Framework**: Echo v4 web framework with OpenAPI-first design
- **Database**: PostgreSQL or SQLite (configurable via `config.yaml`)
- **API**: REST API generated from OpenAPI spec (`internal/api/openapi.yaml`)
- **Structure**: Clean architecture with domain/service/repository patterns
  - `internal/domain/` - Business models and interfaces
  - `internal/infra/` - Infrastructure implementations (DB, repositories)
  - `internal/api/` - HTTP handlers and OpenAPI specification
  - `pkg/` - Shared packages (auth, response utilities)

### Frontend (Vue.js)
- **Framework**: Vue 3 with Composition API + TypeScript
- **UI Library**: Vuetify 3 with Material Design components
- **State Management**: Pinia stores organized by domain
- **Router**: Vue Router with file-based routing via unplugin-vue-router
- **Build Tool**: Vite with TypeScript support
- **API Client**: Auto-generated from OpenAPI spec using openapi-typescript-codegen

## Development Commands

### Backend Development
```bash
# Run development server
go run .

# Build binary
go build -o mosscat.exe

# Generate API handlers from OpenAPI spec
go generate

# Run tests
go test ./...

# Run specific package tests
go test ./internal/infra/repository

# Windows service management
mosscat.exe -service install   # Install as Windows service
mosscat.exe -service uninstall # Remove Windows service
```

### Frontend Development
```bash
cd frontend

# Development server
npm run dev

# Build for production
npm run build

# Type checking
npm run type-check

# Lint and fix
npm run lint

# Generate API client from OpenAPI spec
npm run generate

# Windows users can also use
scripts/generate-client.bat
```

### Testing
```bash
# Install Python test dependencies
pip install -r requirements.txt

# Run integration tests using Tavern (starts server automatically)
pytest tests/
```

## Key Configuration Files

- `config.yaml` - Server configuration (host, port, database DSN, CORS)
- `frontend/vite.config.mts` - Vite build configuration
- `internal/api/openapi.yaml` - API specification (source of truth)

## Database Management

- Migrations in `migrations/` directory using golang-migrate
- Supports both PostgreSQL and SQLite
- Default: SQLite file database (`mosscat.db`)
- Schema includes: users, projects, OSS components/versions, project usages, tags, audit logs

## Code Generation

- Backend API handlers generated from OpenAPI spec using oapi-codegen
- Frontend TypeScript client generated from same OpenAPI spec
- Run `npm run generate` in frontend directory after OpenAPI changes

## Authentication

- JWT-based authentication with configurable roles
- Initial admin user created on first run (password in `admin.initial.password`)
- Role-based access control: Admin, User, Viewer

## Domain Models

Core entities managed by the system:
- **OSS Components**: Software components with metadata and tags
- **OSS Versions**: Specific versions of OSS components
- **Projects**: Internal projects using OSS components
- **Project Usages**: Relationships between projects and OSS versions
- **Tags**: Categorization system for OSS components
- **Users**: System users with role-based permissions
- **Scope Policy**: Rules for OSS usage approval/rejection

## Windows Service Support

Application can run as Windows service via `winservice_windows.go` with service management commands built into the binary.