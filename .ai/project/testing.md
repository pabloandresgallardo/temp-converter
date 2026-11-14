---
Fuente: evidencia observada en código — 2026-11-14
---

# Testing

## Estrategia

TDD estricto: **RED → GREEN → refactor**

1. Escribir el test que falla primero
2. Implementar el mínimo para que pase
3. Refactorizar sin romper tests

## Go (backend)

**Framework:** stdlib `testing` + `testify/assert` y `testify/require`

**Cobertura mínima:** Pendiente de evidencia — se define en primera iteración

**Correr tests:**
```bash
cd backend
GO111MODULE=on go test ./...
GO111MODULE=on go test -v ./internal/domain/...   # solo dominio
```

> Nota: `GO111MODULE=on` requerido — `GO111MODULE=off` está activo en el entorno.

**Tabla tests para conversiones:**
```go
func TestConvert(t *testing.T) {
    cases := []struct{ name string; value float64; from, to Unit; want float64 }{
        {"0°C a F", 0, Celsius, Fahrenheit, 32},
        // ...
    }
    for _, tc := range cases {
        t.Run(tc.name, func(t *testing.T) { ... })
    }
}
```

**Mocks:** interfaces de `port/` se mockean manualmente o con testify/mock.
No usar mockery en esta fase.

## React (frontend)

**Framework:** Vitest + React Testing Library

**Setup:** `src/test-setup.ts` importa `@testing-library/jest-dom`

**Correr tests:**
```bash
cd frontend
node_modules/.bin/vitest run          # run mode (no watch)
node_modules/.bin/vitest              # watch mode
node_modules/.bin/vitest run --coverage  # con cobertura
```

> Nota: `pnpm run test` requiere sesión interactiva (pnpm 12.x). Usar binario directo.

**Principio RTL:** testear comportamiento visible, no implementación interna.
Usar `getByRole`, `getByLabelText` — no `getByTestId` salvo caso extremo.

## Qué se testea

| Capa                      | Estado       | Qué testea                                                        |
|---------------------------|--------------|-------------------------------------------------------------------|
| domain/temperature        | ✅ 10 tests  | Todas las conversiones (tabla tests)                              |
| usecase/ConvertOnly       | ✅ 2 tests   | Execute OK + unidad inválida                                      |
| usecase/ConvertTemperature| ✅ 2 tests   | Execute guarda en repo + no guarda en error                       |
| usecase/GetHistory        | ✅ 2 tests   | Execute vacío + con registros                                     |
| adapter/memory            | ✅ 3 tests   | FindAll vacío, Save+FindAll, orden LIFO                           |
| adapter/http              | ✅ 9 tests   | GET /convert (4), POST /convert (3), GET /history (2)             |
| React features            | ✅ 2 tests   | ConversionForm: submit OK con postConvert, error de API           |
