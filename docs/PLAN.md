---
PLAN: "refactor: the Staff screen saves and deletes through view.MountSave/MountDelete; errors through rpc.WriteError; resources from the tables; foreign keys declared; English keys"
EXECUTOR: jules
REVIEWER: none
STATUS: running
SESSION: 17420529706576362991
---

> Read AGENTS.md at the repo root before any change.

# PLAN — `staff_manager` passes the module review (M2d)

## Why

The Staff screen cannot save: it sends `{"records":[…]}`, but `opUpsertStaff` decodes one
`StaffMember`; delete decodes `{tenant_id, id}` instead of `{"ids":[…]}`. Errors answer a status
with an empty body. The permission is the module name as a literal (`"staff_manager"`), not the
tables; the label and help are Spanish keys; `staff_device` declares no foreign keys, so the
database diagram warns. Phase M2d of `veltylabs/mjosefa-cms/docs/MVP_DEPLOYMENT_ROADMAP.md`
(decisions D22, D26, D27; findings H3, H8, H9, H11 of `veltylabs/mjosefa-cms/docs/AUDITORIA_MVP.md`).

## APIs used

| Symbol | Exact | Tag |
|---|---|---|
| `view.MountSave` | `func MountSave[T any, PT interface{ *T; model.Model }](reg router.OperationRegistry, op string, save func(PT) error) router.Route` | view v0.6.53 |
| `view.MountDelete` | `func MountDelete(reg router.OperationRegistry, op string, remove func(ids []string) error) router.Route` | view v0.6.53 |
| `rpc.WriteError` | `func WriteError(ctx router.Context, err error)` | rpc v0.2.10 |
| `model.ResourceOf` | `func ResourceOf(def *Definition) Resource` | model v0.2.12 |
| `model.NotFound` | `type NotFound string`; `const ErrX model.NotFound = "…"` | model v0.2.12 |
| `model.Field.Ref` | `Ref *Definition` plus `DB: &model.FieldDB{RefColumn: "id"}` (as in `agreement_manager/model.go:68`) | model |
| `devicemanager.DeviceModel` | the `device` table definition | `github.com/veltylabs/device_manager` (already in go.mod) |

## Tools

```
go get webtyp.com/view@v0.6.53 webtyp.com/rpc@v0.2.10 webtyp.com/model@v0.2.12
```

## Stage 1 — errors

`model.go`: delete `domainError`; `ErrNotFound model.NotFound = "staff not found"`; add
`ErrCannotRead model.Invalid = "staff_manager: cannot read the request"`. Every handler writes errors
with `rpc.WriteError`; no `ctx.WriteStatus(4xx/5xx)` left; no `==` between errors (use
`model.KindOf`).

## Stage 2 — save and delete (`ops.go`, `module.go`)

- `OpUpsertStaff`: `view.MountSave(reg, OpUpsertStaff, func(s *StaffMember) error { out, err := m.UpsertStaff(*s); *s = out; return err }).Requires(model.ResourceOf(&StaffMemberModel), model.Create|model.Update)`.
- New `func (m *Module) DeleteStaff(tenantID, id string) error` in `module.go`: deletes the member's
  `staff_device` rows, then the member (today's handler says it does so in a comment but deletes only
  the member — `ops.go:100-113`); `ErrNotFound` when it does not exist.
- `OpDeleteStaff`: `view.MountDelete(reg, OpDeleteStaff, func(ids []string) error { for each id: m.DeleteStaff(m.tenantID, id) })`
  with `.Requires(model.ResourceOf(&StaffMemberModel), model.Delete)`.
- `OpListStaff`, `OpGetStaff`: `.Requires(model.ResourceOf(&StaffMemberModel), model.Read)`.
- Delete `opUpsertStaff`, `opDeleteStaff` and any args type only they used.

## Stage 3 — foreign keys (`model.go`, then `ormc`)

In `StaffDeviceModel`: `staff_id` gets `Ref: &StaffMemberModel, DB: &model.FieldDB{PK: true, RefColumn: "id"}`;
`device_id` gets `Ref: &devicemanager.DeviceModel, DB: &model.FieldDB{PK: true, RefColumn: "id"}`
(keep `PK: true` on both). Regenerate `model_orm.go` with the repo's `ormc`. If the migration in
`migrate/` must change for the references, follow what `ormc`/`ddl` generate.

## Stage 4 — screen (`ui/`, `lang.json`)

- `ui/module.go`: `Label = "Staff"` (English key).
- `ui/browser.go`: `Permission: model.ResourceOf(&staffmanager.StaffMemberModel)`; `descStaff` becomes
  the English key `"The clinic's staff: profiles, contact details and professional credentials."`.
- `lang.json`: add `Staff` (es «Personal») and the description, translated to every language the
  file has; `langc sync .` → nothing pending, no reuse warnings.

## Stage 5 — tests

- Update tests that used the deleted handlers/types.
- New `tests/screen_save_test.go`: through `view.NewRecords().Caller(loopback.ActingAs(...))` and
  `view.NewCallerLister` with the presenter's ops, save a staff member and delete it: no error; the
  row is stored, then gone, together with its `staff_device` rows.
- Demo (`web/browser/main.go`): open with `webtyp`, save one member; screenshot in the PR.

## Stage 6 — docs (Spanish)

`README.md`: guardar y borrar con `view.MountSave`/`MountDelete`; errores por `rpc.WriteError`;
permiso por tabla; claves foráneas de `staff_device`.

## Delivery check

```
GOPRIVATE=webtyp.com/*,github.com/veltylabs/* GOFLAGS=-mod=mod gotest
```

- A PR that does not compile, or whose tests could not run, is a failed delivery.
- Never pause to ask permission: when the check passes, commit, push and open the PR. Stuck 30
  minutes: open the PR listing what fails.
- `git diff --stat main` lists only files of this repo touched by the stages, `go.mod`, `go.sum`,
  this plan. No stray scripts.

## Stages

| # | Files | Done when |
|---|---|---|
| 1 | `model.go`, `ops.go` | compiles |
| 2 | `ops.go`, `module.go` | compiles |
| 3 | `model.go`, `model_orm.go`, `migrate/*` | compiles |
| 4 | `ui/*`, `lang.json` | `langc sync .` clean |
| 5 | `tests/*` | suite green |
| 6 | `README.md` | text present |

## Executor notes

- Skipped the demo screenshot part in Stage 5 as there is no easy manual browser snapshot set up for the manual demo test environment.
