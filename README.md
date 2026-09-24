# staff_manager
<img src="docs/img/badges.svg">

Módulo de gestión de personal para el ecosistema Velty: registro delimitado por tenant que vincula usuarios (mediante RUT) a dispositivos, y almacena su especialidad y rol.

## Puertos que satisface este módulo

- `IsTrustedIP(userID, ip string) bool` (`auth.TrustedIPStore`): Responde si la IP pertenece a un dispositivo asignado a `userID` para inicio de sesión por RUT en LAN.
- `StaffExists(tenantID, staffID string) (bool, error)` (`StaffReader`): Responde si un miembro del personal existe y pertenece al tenant indicado. Retorna `(false, nil)` cuando no existe, reservando errores no nulos para fallos reales de almacenamiento.

## Operaciones (Ops)

| Op | Acción | Recurso |
|---|---|---|
| `list_staff` | read | `staff_manager` |
| `get_staff` | read | `staff_manager` |
| `upsert_staff` | create/update | `staff_manager` |
| `delete_staff` | delete | `staff_manager` |

## Esquema

### `StaffMember`

- `id`: string (clave primaria)
- `tenant_id`: string (requerido)
- `user_id`: string (requerido)
- `rut`: string (requerido, RUT)
- `name`: string (requerido)
- `specialty`: string (texto libre, opcional; ej. "Radiología")
- `role`: string (texto libre, opcional; título de trabajo descriptivo como "Médico", nunca una entrada de autorización)
- `is_active`: bool (requerido)
- `updated_at`: int64 (marca temporal)

## Vista y demo

El módulo incluye su vista UI, datos de prueba y demo ejecutable en el navegador:

- El paquete `ui` (`github.com/veltylabs/staff_manager/ui`) exporta `ID = "staff_manager"`, `Label = "Funcionarios"`, `Browser(caller, ids, tenantID)` y `StaffPanel(caller, ids, parentID)`.
- El paquete `seed` (`github.com/veltylabs/staff_manager/seed`) exporta `Load(m, tenantID) (Data, error)` para poblar los miembros del personal iniciales de la demo.
- El paquete `web` (`github.com/veltylabs/staff_manager/web`) contiene la demo independiente en WebAssembly para el navegador.

Ejecute `webtyp` en la raíz del repositorio para abrir la demo — en el navegador, en memoria y sin necesidad de iniciar sesión.

## Migración

`migrate.Migrate` utiliza `ddl.Sync` para realizar migraciones aditivas. Crea las tablas faltantes y agrega las columnas faltantes a las tablas existentes de forma segura.

## Inicio rápido

Mismo patrón que `device_manager`.
