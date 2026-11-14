---
Fuente: evidencia observada en código — 2026-11-14
---

# Arquitectura

## Patrón

**Backend:** Hexagonal (ports & adapters)

```
internal/
├── domain/temperature/     ← lógica pura, sin dependencias externas
├── port/                   ← interfaces (contratos)
│   └── HistoryRepository
├── usecase/                ← orquestación (depende solo de ports)
│   ├── ConvertOnly         ← implementado: convierte sin side effects
│   ├── ConvertTemperature  ← implementado: convierte + guarda historial
│   └── GetHistory          ← implementado: retorna historial LIFO
└── adapter/
    ├── http/               ← adaptador de entrada (net/http)
    └── memory/             ← adaptador de salida (in-memory store)
```

**Frontend:** Feature-based folders

```
src/
├── features/
│   ├── conversion/         ← formulario, lógica de conversión
│   └── history/            ← lista de conversiones recientes
└── shared/
    └── api/                ← cliente HTTP (fetch al backend)
```

## Stack con versiones

| Componente       | Tecnología                     | Versión  |
|------------------|--------------------------------|----------|
| Backend          | Go                             | 1.23     |
| HTTP             | net/http estándar              | stdlib   |
| Logging          | slog                           | stdlib   |
| Testing Go       | stdlib testing + testify       | 1.9.0    |
| Frontend         | React                          | 18.3.x   |
| Bundler          | Vite                           | 5.4.x    |
| Lenguaje FE      | TypeScript (strict)            | 5.5.x    |
| CSS              | Tailwind CSS                   | 3.4.x    |
| Testing FE       | Vitest + React Testing Library | 2.1.x    |
| Package manager  | pnpm                           | latest   |
| Contenedores     | Docker + docker-compose        | —        |

## Módulo Go

`github.com/gallardo/temp-converter`

## API REST

| Método | Ruta       | Estado       | Descripción                                   |
|--------|------------|--------------|-----------------------------------------------|
| GET    | /convert   | Implementado | Convierte temperatura, retorna resultado JSON (sin side effects) |
| POST   | /convert   | Implementado | Convierte + guarda en historial in-memory     |
| GET    | /history   | Implementado | Retorna historial volátil (LIFO)              |

### Contrato GET /convert

```
GET /convert?from={unit}&to={unit}&value={float64}
```

Unidades válidas: `celsius`, `fahrenheit`, `kelvin`

Respuesta 200: `{"result": float64}`
Respuesta 400: `{"error": string}` (unidad inválida, value no parseable, param faltante)

### Contrato POST /convert

```
POST /convert
Content-Type: application/json
Body: {"value": float64, "from": string, "to": string}
```

Respuesta 200: `{"result": float64}` — convierte y guarda en historial
Respuesta 400: `{"error": string}` (body inválido, unidad inválida)

### Contrato GET /history

```
GET /history
```

Respuesta 200: `[{"value": float64, "from": string, "to": string, "result": float64, "timestamp": int64}]`
Orden: LIFO (más reciente primero). Retorna `[]` cuando no hay registros.

## Historial volátil

Vive in-memory en el servidor Go (`adapter/memory`). Se pierde al reiniciar.
No hay persistencia en DB por diseño.
