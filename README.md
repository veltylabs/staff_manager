# staff_manager
<img src="docs/img/badges.svg">

Staff management module for the Velty ecosystem: tenant-scoped registry linking users (via RUT) to devices, and carrying their specialty and role.

## Ports this module satisfies

- `IsTrustedIP(userID, ip string) bool` (`auth.TrustedIPStore`): Answers whether `ip` belongs to a device assigned to `userID` for LAN RUT login.
- `StaffExists(tenantID, staffID string) (bool, error)` (`StaffReader`): Answers whether a staff member exists and belongs to the given tenant. Returns `(false, nil)` when missing, reserving non-nil errors for genuine storage failures.

## Ops

| Op | Action | Resource |
|---|---|---|
| `list_staff` | read | `staff_manager` |
| `get_staff` | read | `staff_manager` |
| `upsert_staff` | create/update | `staff_manager` |
| `delete_staff` | delete | `staff_manager` |

## Schema

### `StaffMember`

- `id`: string (primary key)
- `tenant_id`: string (required)
- `user_id`: string (required)
- `rut`: string (required, RUT)
- `name`: string (required)
- `specialty`: string (free text, optional; e.g. "Radiología")
- `role`: string (free text, optional; descriptive job title such as "Médico", never an authorization input)
- `is_active`: bool (required)
- `updated_at`: int64 (timestamp)

## Migration

`migrate.Migrate` uses `ddl.Sync` to perform additive migrations. It creates missing tables and adds any missing columns to existing tables safely.

## Quick Start

Same pattern as `device_manager`.
