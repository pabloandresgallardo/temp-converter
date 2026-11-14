# Spec: GET /convert endpoint

Fecha: 2026-11-14
Jira/HDU: Not provided

## Intent

Exponer `GET /convert?from={unit}&to={unit}&value={float}` que retorna el
resultado de la conversión de temperatura. El frontend lo consume y muestra
el resultado al usuario.

## Goals

- Backend: implementar `GET /convert` con los tres query params y respuesta JSON
- Backend: implementar `temperature.Convert` (dominio — actualmente `panic`)
- Backend: implementar `ConvertOnly` use case (hexagonal — sin side effects)
- Frontend: crear API client que llame al endpoint
- Frontend: crear `ConversionForm` en `features/conversion/` con input, selects y resultado

## Non-goals

- Guardar conversión en historial (eso es responsabilidad de `POST /convert`)
- Autenticación
- Validación avanzada más allá de tipos y unidades soportadas

## Anti-goals

- Cambiar o eliminar el stub `POST /convert` existente
- Introducir librerías HTTP externas (mantener `net/http` stdlib)
- Introducir state management global (usar `useState` local)

## Constraints

- Go 1.23 · net/http stdlib · testify 1.9.0
- React 18 · TypeScript strict · Tailwind CSS
- TDD: test escrito antes del código de producción
- Hexagonal: handler no llama domain directamente — pasa por use case

## Approaches

| Approach | Trade-off |
|---|---|
| A) Handler → `ConvertOnly` use case → domain | Sigue hexagonal estricto. Un use case extra para función pura, pero mantiene la convención del proyecto. |
| B) Handler llama domain directamente | Viola hexagonal (adapter→domain sin use case). Más simple pero inconsistente con el patrón del proyecto. |

**Elegido: A** — proyecto de aprendizaje de hexagonal; la consistencia del patrón es el objetivo.

## Contract

### Request

```
GET /convert?from={unit}&to={unit}&value={float64}
```

Valores válidos de `{unit}`: `celsius`, `fahrenheit`, `kelvin`

### Response (200 OK)

```json
{ "result": 32.0 }
```

### Response (400 Bad Request)

```json
{ "error": "invalid unit: 'x'" }
```

Casos de error: unidad desconocida, `value` no parseable como float64.

## Surfaces afectadas

| Archivo | Acción |
|---|---|
| `backend/internal/domain/temperature/converter.go` | Implementar `Convert` |
| `backend/internal/usecase/convert_only.go` | Crear `ConvertOnly` use case |
| `backend/internal/adapter/http/handler.go` | Agregar `GET /convert` route + `handleGetConvert` |
| `backend/cmd/server/main.go` | Wiring: inyectar `ConvertOnly` en handler |
| `frontend/src/shared/api/temperatureApi.ts` | Crear cliente HTTP |
| `frontend/src/features/conversion/ConversionForm.tsx` | Crear componente UI |
| `frontend/src/features/conversion/ConversionForm.test.tsx` | Tests RTL |
| `frontend/src/App.tsx` | Montar `ConversionForm` |

## Tests requeridos (TDD)

### Go

- `TestConvert_CelsiusToFahrenheit` — 0°C → 32°F
- `TestConvert_FahrenheitToCelsius` — 212°F → 100°C
- `TestConvert_CelsiusToKelvin` — 0°C → 273.15K
- `TestConvert_KelvinToCelsius` — 273.15K → 0°C
- `TestConvert_SameUnit` — 100°C → 100°C
- `TestConvert_InvalidUnit` — retorna error
- `TestConvertOnly_Execute` — llama domain, retorna resultado
- `TestHandleGetConvert_OK` — 200 con resultado correcto
- `TestHandleGetConvert_InvalidUnit` — 400
- `TestHandleGetConvert_InvalidValue` — 400

### React

- Render del formulario con inputs
- Submit llama API con los valores correctos
- Resultado se muestra al usuario
- Error de API se muestra al usuario
