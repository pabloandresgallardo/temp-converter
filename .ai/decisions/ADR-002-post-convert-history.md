# ADR-002: POST /convert para conversión con historial

Fecha: 2026-11-14
Jira/HDU: Not provided
Estado: Aceptado

## Contexto

El formulario usaba GET /convert (ConvertOnly, sin side effects). Para registrar
historial volátil, se necesita un endpoint que convierta y persista en memoria.

## Decisión

Implementar POST /convert con body JSON que convierte y guarda en HistoryRepository
in-memory. GET /history lee el repositorio. El formulario frontend pasa a usar POST.

## Alternativas consideradas

| Alternativa | Por qué descartada |
|-------------|-------------------|
| Hacer que GET /convert también guarde | GET no debe tener side effects (RFC 7231) |
| Guardar desde el frontend vía localStorage | Historial no sería compartido entre sesiones ni centralizado |
| Agregar query param `?save=true` a GET | Ambiguo, viola semántica REST, complica caché |

## Consecuencias

- **Positivas:** separación clara de intención (read-only vs. write), historial centralizado en backend, sin duplicación de lógica de conversión
- **Negativas / trade-offs:** formulario frontend debe hacer POST (body JSON) en lugar de GET (query params); conversiones sin historial requieren seguir usando GET /convert
- **Neutrales:** ConvertTemperature reutiliza temperature.Convert del dominio, mismo que ConvertOnly

## Artefactos relacionados

- Spec: `.ai/docs/specs/2026-11-14-post-convert-get-history.md`
- Plan: `.ai/docs/plans/2026-11-14-post-convert-get-history.md`
- ADR previo: `.ai/decisions/ADR-001-get-convert-endpoint.md`
