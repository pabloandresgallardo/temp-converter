---
Fuente: intención de proyecto — 2026-11-14
---

# Convenciones

## Go (backend)

- **Naming:** snake_case para archivos, PascalCase para tipos exportados, camelCase para no exportados
- **Paquetes:** singular, minúscula (`usecase`, `memory`, no `usecases`)
- **Errors:** `fmt.Errorf("usecase: %w", err)` — prefijo del paquete en cada error
- **Interfaces:** definidas en `port/`, implementadas en `adapter/`
- **Tests:** misma carpeta que el código, sufijo `_test.go`
- **Formatter:** `gofmt` (incluido en golangci-lint)
- **Linter:** `golangci-lint run ./...`

## TypeScript / React (frontend)

- **Archivos de componentes:** PascalCase (`ConversionForm.tsx`)
- **Hooks:** camelCase con prefijo `use` (`useConversion.ts`)
- **Servicios / API:** camelCase (`temperatureApi.ts`)
- **Imports:** path alias `@/` para `src/`
- **Formatter:** Prettier (`pnpm format`)
- **Linter:** ESLint (`pnpm lint`)
- **TypeScript:** strict mode — sin `any` explícito

## Commits

Conventional Commits:
- `feat:` nueva funcionalidad
- `fix:` corrección de bug
- `test:` agrega o modifica tests
- `chore:` mantenimiento, deps, config
- `docs:` documentación

## Estructura de feature (frontend)

Cada feature en `src/features/<nombre>/` puede contener:
```
conversion/
├── ConversionForm.tsx
├── ConversionForm.test.tsx
├── useConversion.ts
└── index.ts           ← barrel export
```
