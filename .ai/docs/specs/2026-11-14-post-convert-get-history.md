# Spec: POST /convert + GET /history + historial en frontend

Fecha: 2026-11-14
Jira/HDU: Not provided

## Objetivo

Completar el flujo de historial volátil: cuando el usuario convierte una
temperatura a través del formulario, la conversión se guarda in-memory en el
backend y aparece inmediatamente en la lista de historial del frontend.

## Superficies afectadas

| Superficie | Cambio |
|---|---|
| `adapter/memory/history_repository.go` | Implementar `Save` + `FindAll` con mutex |
| `usecase/convert_temperature.go` | Implementar `Execute`: convertir + guardar |
| `usecase/get_history.go` | Implementar `Execute`: leer del repo |
| `adapter/http/handler.go` | Implementar `handleConvert` + `handleHistory` |
| `cmd/server/main.go` | Cablear repo → use cases → handler |
| `shared/api/temperatureApi.ts` | Agregar `postConvert` (POST con body JSON) |
| `features/conversion/ConversionForm.tsx` | Usar `postConvert`, notificar refresh al padre |
| `features/history/HistoryList.tsx` | Aceptar `refreshKey` prop para re-fetch |
| `App.tsx` | Levantar `refreshKey`, pasarlo a HistoryList |

## Contratos

### POST /convert

```
POST /convert
Content-Type: application/json

Body: {"value": float64, "from": string, "to": string}

200: {"result": float64}
400: {"error": string}   ← param faltante, unidad inválida, value no parseable
```

No usa query params. Body JSON obligatorio.

### GET /history

```
GET /history

200: [{"value": float64, "from": string, "to": string, "result": float64, "timestamp": int64}]
200: []   ← cuando no hay registros (nunca 404)
```

Orden: más reciente primero (LIFO).

## Comportamiento esperado

1. Usuario ingresa valor, selecciona unidades, envía formulario.
2. Frontend hace `POST /convert` con body JSON.
3. Backend convierte, guarda en historial in-memory, responde `{"result": float64}`.
4. Frontend muestra resultado y dispara re-fetch de historial.
5. `GET /history` retorna la lista actualizada.
6. `HistoryList` muestra entradas: `{value} {from} → {result} {to}`.

## Criterios de aceptación

| # | Criterio |
|---|---|
| 1 | `POST /convert` con body válido retorna 200 + `{"result": float64}` |
| 2 | `POST /convert` sin campo requerido retorna 400 |
| 3 | `POST /convert` con unidad inválida retorna 400 |
| 4 | `GET /history` vacío retorna `[]` con 200 |
| 5 | Tras `POST /convert`, `GET /history` incluye el registro |
| 6 | Múltiples conversiones: historial en orden LIFO |
| 7 | Frontend: submit del formulario → resultado visible + historial actualizado |
| 8 | Historial volátil: se borra al reiniciar el backend |

## Restricciones

- `HistoryRepository` usa `sync.RWMutex` (struct ya tiene el campo).
- `ConvertTemperature` delega conversión a `temperature.Convert` (mismo dominio).
- `Timestamp` = Unix epoch en segundos (`time.Now().Unix()`).
- No agregar dependencias externas.
- No persistencia en disco ni DB.
- `GET /convert` no se modifica — sigue siendo ConvertOnly sin side effects.
