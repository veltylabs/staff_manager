# staff_manager
<img src="docs/img/badges.svg">

Staff management module for the Velty ecosystem: tenant-scoped registry linking users (via RUT) to devices. Implements `auth.TrustedIPStore` (`IsTrustedIP`) for LAN RUT login.

## Ops

| Op | Action | Resource |
|---|---|---|
| `list_staff` | read | `staff_manager` |
| `get_staff` | read | `staff_manager` |
| `upsert_staff` | create/update | `staff_manager` |
| `delete_staff` | delete | `staff_manager` |

## Quick Start

Same pattern as `device_manager`.

