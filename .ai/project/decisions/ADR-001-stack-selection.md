# ADR-001: Selección de stack tecnológico

Fecha: 2026-11-14
Jira/HDU: Proyecto de aprendizaje / prueba
Estado: Aceptado

## Contexto

App fullstack de conversión de temperaturas (Celsius, Fahrenheit, Kelvin)
con historial volátil in-memory. Propósito: aprendizaje de arquitectura
hexagonal en Go y feature-based React con TDD estricto.

## Decisión

**Backend:** Go 1.23 · net/http estándar · slog · testify 1.9.0
**Patrón:** Hexagonal (ports & adapters)
**Frontend:** React 18.3 · Vite 5.4 · TypeScript 5.5 (strict) · Tailwind 3.4
**Testing FE:** Vitest 2.1 + React Testing Library 16
**State:** useState + Context API
**Package manager:** pnpm
**Linting:** golangci-lint (Go) · ESLint + Prettier (React)
**Infra:** Docker + docker-compose
**Commits:** Conventional Commits

## Alternativas consideradas

| Alternativa | Por qué descartada |
|---|---|
| Go + Gin/Chi router | net/http 1.22+ cubre el caso; añade dependencia sin ventaja real para esta escala |
| Go + Standard architecture (handlers→services→repos) | Usuario eligió hexagonal por objetivo de aprendizaje del patrón |
| React 19 | Algunas librerías (RTL, etc.) aún no 100% compatibles; React 18 es estable |
| Next.js en lugar de Vite | Overkill para SPA sin SSR; agrega complejidad innecesaria |
| Jest en lugar de Vitest | Vitest es nativo de Vite, más rápido, sin config extra |
| Zustand / Redux Toolkit | Context API suficiente para historial global en app de esta escala |
| State local React (sin backend) | Usuario requirió historial in-memory en servidor Go |
| zap / logrus para logging | slog stdlib evita dependencia externa; logrus en modo mantenimiento |
| GitHub Actions CI/CD | No planeado en esta etapa |

## Consecuencias

**Positivas:**
- Dominio (`domain/`) completamente aislado de infraestructura — testeable sin mocks
- net/http sin dependencias externas — aprende Go puro
- Vite + Vitest: dev server rápido, testing integrado sin config extra
- Tailwind elimina necesidad de CSS custom para UI de esta escala

**Negativas / trade-offs:**
- Hexagonal agrega boilerplate (ports, adapters) no estrictamente necesario para este dominio simple — costo aceptado por objetivo de aprendizaje
- Historial volátil (in-memory) se pierde al reiniciar — por diseño, no bug
- pnpm requiere instalación previa si no está en el sistema

## Artefactos relacionados

- Fundaciones: `.ai/project/architecture.md`
- Diagrama de arquitectura: `.ai/project/diagrams/architecture-overview.md`
- Flujo de datos: `.ai/project/diagrams/data-flow.md`
