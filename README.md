# temp-converter
AI-ENGINEERING-MANAGED

## Overview

App fullstack que convierte temperaturas entre Celsius, Fahrenheit y Kelvin,
con historial volátil de conversiones recientes. El historial vive in-memory
en el servidor Go y se pierde al reiniciar.

| Capa      | Tecnología                          |
|-----------|-------------------------------------|
| Frontend  | React 18 + Vite + TypeScript        |
| CSS       | Tailwind CSS                        |
| Backend   | Go 1.23 · net/http estándar         |
| Patrón    | Hexagonal (ports & adapters)        |
| Testing   | testify (Go) · Vitest + RTL (React) |
| Infra     | Docker + docker-compose             |

## Getting started

```bash
# Backend (requiere Go 1.23+)
cd backend && go run ./cmd/server

# Frontend (requiere pnpm)
cd frontend && pnpm install && pnpm dev

# Con Docker
docker-compose up
```

## Architecture

```
temp-converter/
├── backend/      # Go 1.23 — hexagonal (ports & adapters)
├── frontend/     # React 18 + Vite — feature-based folders
├── docker-compose.yml
└── .ai/project/  # arquitectura, decisiones, diagramas
```

Ver `.ai/project/architecture.md` para detalle de capas y módulo Go.

## Development

- Linting Go: `cd backend && golangci-lint run ./...`
- Linting React: `cd frontend && pnpm lint`
- Format React: `cd frontend && pnpm format`
- Commits: Conventional Commits (`feat:`, `fix:`, `chore:`, `test:`)

## Testing

- Backend: `cd backend && go test ./...`
- Frontend: `cd frontend && pnpm test`
- Estrategia: TDD estricto (RED → GREEN → refactor)

Ver `.ai/project/testing.md` para framework, cobertura y qué testear por capa.

## Documentation

- Arquitectura: `.ai/project/architecture.md`
- Convenciones: `.ai/project/conventions.md`
- Testing: `.ai/project/testing.md`
- Guardrails: `.ai/project/guardrails.md`
- ADR-001: `.ai/project/decisions/ADR-001-stack-selection.md`
- Diagramas: `.ai/project/diagrams/`

## Integrations

Not detected

## Guardrails

Ver `.ai/project/guardrails.md`. Restricciones críticas:

- `golangci-lint` y `pnpm lint` deben pasar sin errores antes de merge
- No hardcodear secrets ni URLs de producción en código
- Dominio (`internal/domain/`) sin imports externos al módulo Go
- TypeScript strict: sin `any` explícito
