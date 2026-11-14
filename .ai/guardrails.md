---
Fuente: intención de proyecto — 2026-11-14
---

# Guardrails

## Blocking (bloquean merge/entrega)

- `golangci-lint run ./...` debe pasar sin errores (backend)
- `pnpm lint` debe pasar sin warnings (frontend)
- No hardcodear secrets, tokens ni URLs de producción en código
- No commitear `.env` ni archivos con credenciales
- TDD: no código de producción sin test previo (excepto scaffolding inicial)
- TypeScript strict: no `any` explícito, no ignorar errores de tipo

## Required (requerido por convención)

- Conventional Commits en todos los commits
- Interfaces en `port/`, implementaciones en `adapter/`
- Dominio (`domain/`) sin imports de paquetes externos al módulo
- Tests junto al código (misma carpeta, sufijo `_test.go` / `.test.tsx`)
- Barrel exports (`index.ts`) en cada feature de React

## Guidance (preferencias del developer)

- Preferir `net/http` estándar sobre frameworks en Go
- Preferir `slog` para logging sobre librerías externas
- Tailwind utilities directamente en JSX — no crear clases CSS custom salvo necesidad clara
- Context API para estado global — no introducir Zustand/Redux sin consenso
- Commits pequeños y frecuentes, un concepto por commit
