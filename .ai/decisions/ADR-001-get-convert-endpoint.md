# ADR-001: GET /convert con ConvertOnly use case

Fecha: 2026-11-14
Jira/HDU: Not provided
Estado: Aceptado

## Contexto

Se necesita exponer la conversión de temperatura como endpoint HTTP consumible
desde React. La arquitectura del proyecto es hexagonal (ports & adapters).

## Decisión

Usar `GET /convert?from=&to=&value=` con un use case `ConvertOnly` sin ports
ni side effects. El handler (adapter) llama al use case, no al domain directamente.

## Alternativas consideradas

| Alternativa | Por qué descartada |
|---|---|
| Handler llama `temperature.Convert` directamente | Viola hexagonal: adapter no debe llamar domain sin pasar por use case |
| Reusar `ConvertTemperature` (con historial) | Ese use case guarda en historial; GET debe ser sin side effects |
| POST /convert con body JSON | Query params son más simples para conversión stateless; GET es semánticamente correcto para operación idempotente |

## Consecuencias

- **Positivas:** endpoint semánticamente correcto (GET idempotente), use case puro sin ports (fácil de testear sin mocks), patrón hexagonal respetado
- **Negativas / trade-offs:** un use case extra para función que podría llamarse directamente; overhead mínimo justificado por coherencia del patrón
- **Neutrales:** `POST /convert` stub existente se mantiene para futura integración con historial

## Artefactos relacionados

- Spec: `.ai/docs/specs/2026-11-14-get-convert-endpoint.md`
- Plan: `.ai/docs/plans/2026-11-14-get-convert-endpoint.md`
