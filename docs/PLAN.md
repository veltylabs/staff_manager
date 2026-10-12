---
PLAN: "feat: StaffCard — nombre, especialidad y estado de un profesional, sin datos personales"
EXECUTOR: jules
REVIEWER: none
STATUS: review
SESSION: 4804385620053531228
PR: https://github.com/veltylabs/staff_manager/pull/5
---

> Este plan se despacha con el flujo CodeJob. Ver skill: agents-workflow.

> **Pieza S** de
> [`mjosefa-cms/docs/RESERVA_Y_PAGOS_MASTER_PLAN.md`](https://github.com/veltylabs/mjosefa-cms/blob/main/docs/RESERVA_Y_PAGOS_MASTER_PLAN.md).
> Compuerta de la pieza B (`appointment_booking`).

# Plan — `staff_manager`: método `StaffCard`

## 0. Contexto

La reserva online (pública) muestra el nombre y la especialidad de cada profesional.
`appointment_booking` declara su propio lector `StaffReader` y no puede importar este módulo
(regla: un paquete raíz no importa a otro). Hoy `StaffReader` solo tiene `StaffExists`. Hace
falta un método con **tipos primitivos** en la firma, para que `*Module` lo satisfaga
estructuralmente, que devuelva solo lo que se puede mostrar al público.

## Design gate

1. **Prior art:** el propio `StaffExists` (puerto estrecho satisfecho estructuralmente);
   FHIR `PractitionerRole` público (nombre + especialidad, sin identificadores personales).
2. **Prueba del nombre:** `StaffCard(tenantID, staffID)` "la tarjeta del profesional".
3. **Libro de complejidad:** +1 método exportado. Nada más cambia.
4. **Dónde vive:** aquí; el dueño del personal.
5. **Qué borra:** nada.

## 1. Reglas

- Sin `strings`/`strconv`/`errors`/`fmt` de la stdlib (usar `webtyp.com/fmt`); sin `map`; sin
  `reflect`. Comentarios en español.
- **Todas las pruebas en `tests/`** (paquete `tests`, solo API exportada). `gotest`
  (prerrequisito: `go install webtyp.com/devflow/cmd/gotest@latest`).
- **TDD:** primero la firma con `panic("sin implementar")` y las pruebas de §3 → `gotest` →
  confirmar que fallan → commit `test: StaffCard en rojo` → implementar → verde → commit
  `feat: StaffCard`.
- **Cobertura:** `StaffCard` en 100.0%:
  `go test -coverpkg=github.com/veltylabs/staff_manager -coverprofile=/tmp/c.out ./tests/... && go tool cover -func=/tmp/c.out | grep StaffCard`.
- Diagrama: [`docs/diagrams/staff_card.md`](diagrams/staff_card.md) (ya escrito).

## 2. Implementación (`module.go`, junto a `StaffExists`)

```go
// StaffCard devuelve lo que se puede mostrar al público de un profesional:
// nombre, especialidad y si está activo. Nunca RUT, usuario, dispositivos ni
// contacto. Satisface StaffReader.StaffCard de appointment_booking
// estructuralmente. Ausente u otro tenant ⇒ ErrNotFound.
func (m *Module) StaffCard(tenantID, staffID string) (name, specialty string, active bool, err error)
```

Diagrama §1: vacío ⇒ `ErrNotFound`; `GetStaff` ⇒ `ErrNotFound` tal cual; otro error tal cual;
encontrado ⇒ `s.Name, s.Specialty, s.IsActive, nil`.

## 3. Pruebas — trazabilidad (`tests/staff_card_test.go`)

| Rama | Prueba |
|---|---|
| `K1-sí` | `TestStaffCard_EmptyArgsNotFound` (tenant vacío, id vacío) |
| `K2-encontrado` | `TestStaffCard_ReturnsNameSpecialtyActive` (activo e inactivo) |
| `K2-ErrNotFound` | `TestStaffCard_UnknownOrOtherTenantNotFound` |
| `K2-otro error` | `TestStaffCard_StorageErrorPropagates` (un `storage.Conn` falso declarado en `tests/` que envuelve `storage/mem` y falla en `Query`) |

## 4. Documentación

`README.md`: una línea en la sección de lectores para otros módulos: `StaffCard` y para qué.

## 5. Criterios de aceptación

- `gotest` en verde; `StaffCard` en 100.0%; cada prueba de §3 con su nombre exacto.
- Historial: commit `test:` en rojo antes del `feat:`.
- `git diff --stat main` lista exactamente `module.go`, `tests/staff_card_test.go`, `README.md`;
  nada borrado.
