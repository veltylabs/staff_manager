# Bugs de UI y mejoras — `staff_manager`

> Auditoría en vivo sobre `#staff_manager` (2026-11-19).

## Bugs encontrados

| # | Sev | Pantalla | Qué se ve en pantalla | Causa raíz verificada |
|---|---|---|---|---|
| SM-1 | 🔴 | `#staff_manager` | **«Especialidad» y «Rol» son campos de texto libre** (`"Odontología"`, `"cirujano dental"`) en lugar de selectores vinculados al catálogo y a los roles del sistema. | `StaffMemberModel` (`staff_manager/definition.go:23-24`) declara `{Name: "specialty", Type: input.Text()}` y `{Name: "role", Type: input.Text()}` violando `ui-consistency` («No free-text codes»). |
| SM-2 | 🟡 | `#staff_manager` | **Descalce real con `service_catalog` que rompe la reserva de horas:** el personal tiene escrito `"Odontología"` (texto libre) mientras el catálogo usa `"Dental (Odontología + Ortodoncia)"`. | Al no existir selector contra `service_catalog`, los nombres divergen y el personal queda desconectado de especialidades y roles reales. |
| SM-3 | 🟡 | Menú | **Colisión de nombre «Personal»** con `#personal` de `appointment_booking`. | Ambos módulos usan `"Personal"` en la barra lateral de `mjosefa-cms`. |

## Sugerencias de mejora justificadas

1. **Convertir `Especialidad` y `Rol` en selectores reales alimentados por sus módulos dueños (o eliminar duplicidad de Rol si lo gestiona `access`):** inyectar las especialidades de `service_catalog` en `staff_manager` vía interfaz/reader (como ya hace `appointment_booking`). *Justificación:* evita errores de tipeo (`"Odontología"` vs `"odontologia"` vs `"Dental"`) y garantiza que todo profesional elegible aparezca correctamente en agenda y box.
2. **Mostrar en la fila de la lista si el funcionario está inactivo y su especialidad (badge limpio en lugar de solo el RUT):** *Justificación:* facilita distinguir de un vistazo médicos activos de recepcionistas o ex-funcionarios (57 registros).
