# `StaffCard` — lo mínimo y público de un profesional

Cada rombo lleva un id; [`docs/PLAN.md`](../PLAN.md) §3 asigna a cada rama su prueba.

```mermaid
flowchart TD
    A[StaffCard tenant, id] --> K1{K1: ¿tenant o id vacío?}
    K1 -- sí --> NF0[ErrNotFound]
    K1 -- no --> G[GetStaff]
    G --> K2{K2: ¿resultado?}
    K2 -- encontrado --> OK[nombre, especialidad, activo]
    K2 -- ErrNotFound --> NF[ErrNotFound]
    K2 -- otro error --> E[devolver el error]
```

Nunca devuelve RUT, usuario, dispositivos ni contacto: lo consume una reserva online pública.
