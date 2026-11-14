# GET /convert Endpoint — Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Exponer `GET /convert?from=&to=&value=` en Go y consumirlo desde React con un formulario de conversión.

**Architecture:** Hexagonal: handler (adapter) → ConvertOnly use case → temperature.Convert (domain). El use case es una función pura sin ports ni side effects. El frontend llama al endpoint con fetch y muestra el resultado inline.

**Tech Stack:** Go 1.23 · net/http stdlib · testify 1.9.0 · React 18 · TypeScript strict · Tailwind CSS · Vitest + RTL

**Spec:** `.ai/docs/specs/2026-11-14-get-convert-endpoint.md`

## Global Constraints

- Go 1.23 — no introducir versiones de sintaxis posteriores
- net/http stdlib — sin frameworks HTTP externos
- testify v1.9.0 — `assert` y `require` únicamente
- TypeScript strict — sin `any` explícito
- Tailwind CSS — sin clases CSS custom
- TDD: test escrito y fallando antes del código de producción
- Hexagonal: handler no llama domain directamente

## Review Focus

1. **Query params ausentes** (`from`, `to`, o `value` faltantes): debe retornar 400 con mensaje descriptivo — cubierto en Task 3, Step 1.
2. **Kelvin negativo** (`value=-1&from=kelvin&to=celsius`): temperaturas por debajo de 0K son físicamente inválidas; la spec no define si debe aceptarse o rechazarse — plan lo acepta (no es una unidad inválida, solo un valor físicamente imposible). Cubierto en Task 1 implícitamente: `Convert` no valida rangos de valor.
3. **Valor entero en query string** (`value=100`): `strconv.ParseFloat("100", 64)` lo parsea correctamente — cubierto en Task 3, Step 1.
4. **Unidad con casing diferente** (`from=Celsius`): la spec define valores lowercase; mayúsculas deben retornar 400 — cubierto en Task 1 `TestConvert_InvalidUnit`.
5. **Same-unit conversion** (`from=celsius&to=celsius`): debe retornar el valor sin cambios — cubierto en Task 1 `TestConvert_SameUnit`.

---

### Task 1: Domain — temperature.Convert

**Files:**
- Modify: `backend/internal/domain/temperature/converter.go`
- Create: `backend/internal/domain/temperature/converter_test.go`

**Interfaces:**
- Produces: `Convert(value float64, from, to Unit) (float64, error)` — retorna error si unit es inválida, resultado float64 si válida

- [ ] **Step 1: Escribir tests fallidos**

```go
// backend/internal/domain/temperature/converter_test.go
package temperature_test

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "github.com/gallardo/temp-converter/internal/domain/temperature"
)

func TestConvert(t *testing.T) {
    cases := []struct {
        name  string
        value float64
        from  temperature.Unit
        to    temperature.Unit
        want  float64
    }{
        {"0C→F",   0,      temperature.Celsius,    temperature.Fahrenheit, 32},
        {"100C→F", 100,    temperature.Celsius,    temperature.Fahrenheit, 212},
        {"212F→C", 212,    temperature.Fahrenheit, temperature.Celsius,    100},
        {"0C→K",   0,      temperature.Celsius,    temperature.Kelvin,     273.15},
        {"273.15K→C", 273.15, temperature.Kelvin,  temperature.Celsius,    0},
        {"32F→K",  32,     temperature.Fahrenheit, temperature.Kelvin,     273.15},
        {"273.15K→F", 273.15, temperature.Kelvin,  temperature.Fahrenheit, 32},
        {"SameUnit", 100,  temperature.Celsius,    temperature.Celsius,    100},
    }
    for _, tc := range cases {
        t.Run(tc.name, func(t *testing.T) {
            got, err := temperature.Convert(tc.value, tc.from, tc.to)
            require.NoError(t, err)
            assert.InDelta(t, tc.want, got, 0.001)
        })
    }
}

func TestConvert_InvalidUnit(t *testing.T) {
    _, err := temperature.Convert(100, "rankine", temperature.Celsius)
    assert.ErrorContains(t, err, "invalid unit")
}
```

- [ ] **Step 2: Correr tests — verificar que fallan**

```bash
cd backend && go test ./internal/domain/temperature/... -v
```
Expected: FAIL — `panic: not implemented`

- [ ] **Step 3: Implementar `Convert` en `converter.go`**

Fórmulas exactas:
- C→F: `(v * 9/5) + 32`
- F→C: `(v - 32) * 5/9`
- C→K: `v + 273.15`
- K→C: `v - 273.15`
- F→K: `(v - 32) * 5/9 + 273.15`
- K→F: `(v - 273.15) * 9/5 + 32`
- Same unit: retornar `v, nil`
- Unidad desconocida: retornar `0, fmt.Errorf("invalid unit: %q", from)` (o `to` si `to` es inválida)

- [ ] **Step 4: Correr tests — verificar PASS**

```bash
cd backend && go test ./internal/domain/temperature/... -v
```
Expected: PASS todos los casos

- [ ] **Step 5: Commit**

```bash
git add backend/internal/domain/temperature/
git commit -m "feat(domain): implement temperature Convert function"
```

---

### Task 2: ConvertOnly use case

**Files:**
- Create: `backend/internal/usecase/convert_only.go`
- Create: `backend/internal/usecase/convert_only_test.go`

**Interfaces:**
- Consumes: `temperature.Convert(float64, Unit, Unit) (float64, error)` de Task 1
- Produces:
  - `type ConvertOnlyRequest struct { Value float64; From temperature.Unit; To temperature.Unit }`
  - `type ConvertOnlyResponse struct { Result float64 }`
  - `type ConvertOnly struct{}`
  - `func NewConvertOnly() *ConvertOnly`
  - `func (uc *ConvertOnly) Execute(req ConvertOnlyRequest) (ConvertOnlyResponse, error)`

- [ ] **Step 1: Escribir test fallido**

```go
// backend/internal/usecase/convert_only_test.go
package usecase_test

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "github.com/gallardo/temp-converter/internal/domain/temperature"
    "github.com/gallardo/temp-converter/internal/usecase"
)

func TestConvertOnly_Execute(t *testing.T) {
    uc := usecase.NewConvertOnly()
    resp, err := uc.Execute(usecase.ConvertOnlyRequest{
        Value: 0,
        From:  temperature.Celsius,
        To:    temperature.Fahrenheit,
    })
    require.NoError(t, err)
    assert.InDelta(t, 32.0, resp.Result, 0.001)
}

func TestConvertOnly_Execute_InvalidUnit(t *testing.T) {
    uc := usecase.NewConvertOnly()
    _, err := uc.Execute(usecase.ConvertOnlyRequest{
        Value: 0,
        From:  "rankine",
        To:    temperature.Celsius,
    })
    assert.Error(t, err)
}
```

- [ ] **Step 2: Correr test — verificar que falla**

```bash
cd backend && go test ./internal/usecase/... -run TestConvertOnly -v
```
Expected: FAIL — tipo no existe

- [ ] **Step 3: Implementar `ConvertOnly` en `convert_only.go`**

Signatures exactas de Interfaces → Produces. `Execute` delega a `temperature.Convert`.

- [ ] **Step 4: Correr tests — verificar PASS**

```bash
cd backend && go test ./internal/usecase/... -v
```
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add backend/internal/usecase/convert_only.go backend/internal/usecase/convert_only_test.go
git commit -m "feat(usecase): add ConvertOnly use case"
```

---

### Task 3: HTTP adapter — GET /convert

**Files:**
- Modify: `backend/internal/adapter/http/handler.go`
- Create: `backend/internal/adapter/http/handler_test.go`
- Modify: `backend/cmd/server/main.go`

**Interfaces:**
- Consumes: `usecase.ConvertOnly` y `usecase.ConvertOnlyRequest`/`ConvertOnlyResponse` de Task 2
- Produces:
  - `GET /convert?from={unit}&to={unit}&value={float}` → `200 {"result": float64}` o `400 {"error": string}`
  - `NewHandler` signature nueva: `func NewHandler(convert *usecase.ConvertTemperature, history *usecase.GetHistory, convertOnly *usecase.ConvertOnly) *Handler`

- [ ] **Step 1: Escribir tests fallidos**

```go
// backend/internal/adapter/http/handler_test.go
package http_test

import (
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    httpadapter "github.com/gallardo/temp-converter/internal/adapter/http"
    "github.com/gallardo/temp-converter/internal/usecase"
)

func newTestHandler() *httpadapter.Handler {
    return httpadapter.NewHandler(nil, nil, usecase.NewConvertOnly())
}

func TestHandleGetConvert_OK(t *testing.T) {
    h := newTestHandler()
    mux := http.NewServeMux()
    h.Register(mux)

    req := httptest.NewRequest("GET", "/convert?from=celsius&to=fahrenheit&value=0", nil)
    rec := httptest.NewRecorder()
    mux.ServeHTTP(rec, req)

    assert.Equal(t, http.StatusOK, rec.Code)
    var resp map[string]float64
    require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))
    assert.InDelta(t, 32.0, resp["result"], 0.001)
}

func TestHandleGetConvert_MissingParam(t *testing.T) {
    h := newTestHandler()
    mux := http.NewServeMux()
    h.Register(mux)

    req := httptest.NewRequest("GET", "/convert?from=celsius&to=fahrenheit", nil) // missing value
    rec := httptest.NewRecorder()
    mux.ServeHTTP(rec, req)

    assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestHandleGetConvert_InvalidUnit(t *testing.T) {
    h := newTestHandler()
    mux := http.NewServeMux()
    h.Register(mux)

    req := httptest.NewRequest("GET", "/convert?from=rankine&to=celsius&value=100", nil)
    rec := httptest.NewRecorder()
    mux.ServeHTTP(rec, req)

    assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestHandleGetConvert_InvalidValue(t *testing.T) {
    h := newTestHandler()
    mux := http.NewServeMux()
    h.Register(mux)

    req := httptest.NewRequest("GET", "/convert?from=celsius&to=fahrenheit&value=abc", nil)
    rec := httptest.NewRecorder()
    mux.ServeHTTP(rec, req)

    assert.Equal(t, http.StatusBadRequest, rec.Code)
}
```

- [ ] **Step 2: Correr tests — verificar que fallan**

```bash
cd backend && go test ./internal/adapter/http/... -v
```
Expected: FAIL — compilación o "not implemented"

- [ ] **Step 3: Modificar `handler.go`**

Agregar campo `convertOnlyUC *usecase.ConvertOnly` al struct `Handler`.
Actualizar `NewHandler` a tres params (mantener `convert` y `history` como los existentes, agregar `convertOnly *usecase.ConvertOnly`).
En `Register`, agregar `mux.HandleFunc("GET /convert", h.handleGetConvert)`.
Implementar `handleGetConvert`:
- Leer `from`, `to`, `value` de `r.URL.Query()`; si alguno está vacío → 400 `{"error":"missing param: <name>"}`
- `strconv.ParseFloat(value, 64)` — si falla → 400 `{"error":"invalid value"}`
- Llamar `h.convertOnlyUC.Execute(...)` — si error → 400 `{"error": err.Error()}`
- 200 con `json.NewEncoder(w).Encode(map[string]float64{"result": resp.Result})`
- Setear `Content-Type: application/json` en todos los casos

- [ ] **Step 4: Correr tests — verificar PASS**

```bash
cd backend && go test ./internal/adapter/http/... -v
```
Expected: PASS todos

- [ ] **Step 5: Actualizar wiring en `main.go`**

Instanciar `usecase.NewConvertOnly()` y pasarlo como tercer argumento a `NewHandler`.

- [ ] **Step 6: Verificar compilación**

```bash
cd backend && go build ./...
```
Expected: sin errores

- [ ] **Step 7: Commit**

```bash
git add backend/internal/adapter/http/ backend/cmd/server/main.go
git commit -m "feat(api): add GET /convert endpoint"
```

---

### Task 4: Frontend — API client + ConversionForm + wiring

**Files:**
- Create: `frontend/src/shared/api/temperatureApi.ts`
- Create: `frontend/src/features/conversion/ConversionForm.tsx`
- Create: `frontend/src/features/conversion/ConversionForm.test.tsx`
- Create: `frontend/src/features/conversion/index.ts`
- Modify: `frontend/src/App.tsx`

**Interfaces:**
- Consumes: `GET /convert?from=&to=&value=` → `{result: number}` de Task 3
- Produces:
  - `convertTemperature(value: number, from: string, to: string): Promise<{ result: number }>` en `temperatureApi.ts`
  - `<ConversionForm />` — componente sin props, estado local

- [ ] **Step 1: Escribir test fallido**

```tsx
// frontend/src/features/conversion/ConversionForm.test.tsx
import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { vi } from 'vitest'
import { ConversionForm } from './ConversionForm'
import * as api from '@/shared/api/temperatureApi'

vi.mock('@/shared/api/temperatureApi')

describe('ConversionForm', () => {
  it('muestra el resultado después de submit', async () => {
    vi.mocked(api.convertTemperature).mockResolvedValue({ result: 32 })
    render(<ConversionForm />)
    await userEvent.type(screen.getByLabelText(/valor/i), '0')
    fireEvent.submit(screen.getByRole('form'))
    await waitFor(() => expect(screen.getByText(/32/)).toBeInTheDocument())
  })

  it('muestra error cuando la API falla', async () => {
    vi.mocked(api.convertTemperature).mockRejectedValue(new Error('invalid unit'))
    render(<ConversionForm />)
    await userEvent.type(screen.getByLabelText(/valor/i), '0')
    fireEvent.submit(screen.getByRole('form'))
    await waitFor(() => expect(screen.getByRole('alert')).toBeInTheDocument())
  })
})
```

- [ ] **Step 2: Correr test — verificar que falla**

```bash
cd frontend && pnpm test --run ConversionForm
```
Expected: FAIL — módulo no existe

- [ ] **Step 3: Crear `temperatureApi.ts`**

Signature: `export async function convertTemperature(value: number, from: string, to: string): Promise<{ result: number }>`
Usa `fetch` con la URL `/convert?from=${from}&to=${to}&value=${value}`.
Si `!response.ok` → throw Error con texto del body.

- [ ] **Step 4: Crear `ConversionForm.tsx`**

Estado local: `value: string`, `from: string` (default `"celsius"`), `to: string` (default `"fahrenheit"`), `result: number | null`, `error: string | null`.
Render: `<form aria-label="form">` con `<label htmlFor="value">Valor</label>`, `<input id="value">`, dos `<select>` para `from`/`to` con opciones `celsius`, `fahrenheit`, `kelvin`, botón submit.
En submit: llama `convertTemperature`, setea `result` o `error`.
Muestra `result` como texto y `error` en `<div role="alert">`.

- [ ] **Step 5: Crear `index.ts`**

```ts
export { ConversionForm } from './ConversionForm'
```

- [ ] **Step 6: Correr tests — verificar PASS**

```bash
cd frontend && pnpm test --run ConversionForm
```
Expected: PASS

- [ ] **Step 7: Montar en `App.tsx`**

```tsx
import { ConversionForm } from '@/features/conversion'
// reemplazar el TODO comentario con <ConversionForm />
```

- [ ] **Step 8: Verificar build**

```bash
cd frontend && pnpm build
```
Expected: sin errores TypeScript ni de build

- [ ] **Step 9: Commit**

```bash
git add frontend/src/shared/api/ frontend/src/features/conversion/ frontend/src/App.tsx
git commit -m "feat(frontend): add ConversionForm and API client"
```
