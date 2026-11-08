---
PLAN: "refactor!: staff_manager on the closed crud contract — server tenant, transactional batch, single declarative validation, view.go browser-only"
EXECUTOR: jules
REVIEWER: none
STATUS: running
SESSION: 15975789775639801906
---

> Read AGENTS.md at the repo root before any change.
> **Finish by yourself:** when the delivery check passes, commit, push and open the pull request. Never ask whether to
> proceed or finalize — nobody answers, and a session that waits ends without your work. Commit after every stage.
> Never run `git reset`, `git restore`, `git checkout -- <file>` or `git stash` on your own work.
> Never edit `model_orm.go` or `model_crud.go`: a build error there means the hand-written code is wrong.

# Plan: `staff_manager` on `crud.Registry`

Context: master `webtyp/docs/ARCHITECTURE_MASTER_PLAN.md` D2–D16 (phases A4.2 + L5). Follows the template in `patient_directory`.
`ops.go` currently mounts raw operations (`OpListStaff`, `OpGetStaff`, `OpUpsertStaff`, `OpDeleteStaff`).
After this plan the module mounts operations through generated `MountStaffMember*` helpers on `crud.Registry`,
tenant comes from `Deps`, validation of the schema lives in one place (`crud`), and `seed.Demo` provides the demo.

## Tools and local dependencies

```bash
go install webtyp.com/devflow/cmd/gotest@latest
go install webtyp.com/lang/cmd/langc@latest
go install webtyp.com/ormc/cmd/ormc@v0.1.36
```
```bash
go get webtyp.com/model@v0.2.14 webtyp.com/router@v0.4.19 webtyp.com/crud@v0.0.3 webtyp.com/rpc@v0.2.17 webtyp.com/mcp@v0.3.23 \
  webtyp.com/dbsync@v0.3.47 webtyp.com/view@latest webtyp.com/appkit@latest
```
`view@latest` and `appkit@latest` contain `view.Demo` and `appkit` without `layout` (plans L1, L2).
Remove every `replace` from `go.mod`; `go mod tidy`; run `ormc` (generates `model_crud.go`).
After the delivery check, `go.mod` does not list `webtyp.com/layout`.

## APIs used

| Symbol | Signature |
|---|---|
| generated `MountStaffMemberRead` | `func MountStaffMemberRead(reg crud.Registry, newFilter func() model.Model, read func(c crud.Caller, args crud.ReadArgs, filter model.Model) (StaffMemberList, error)) crud.Route` |
| generated `MountStaffMemberCreate` | `func MountStaffMemberCreate(reg crud.Registry, create func(c crud.Caller, list StaffMemberList) error) crud.Route` |
| generated `MountStaffMemberUpdate` | `func MountStaffMemberUpdate(reg crud.Registry, update func(c crud.Caller, ids, fields []string, rec StaffMember) error) crud.Route` |
| generated `MountStaffMemberDelete` | `func MountStaffMemberDelete(reg crud.Registry, deleteFn func(c crud.Caller, ids []string) error) crud.Route` |
| `crud.ReadArgs` | `struct{ ID, Search, After string; Limit int64; Filter model.Model }` |
| `crud.Loopback` / `crud.LoopbackActingAs` | `func Loopback(mods ...crud.Module) router.Caller` / `func LoopbackActingAs(userID string, mods ...crud.Module) router.Caller` |
| `crud.Op` | `func Op(m model.ModuleNaming, action model.Action) string` → `staff_member_r`, `staff_member_c`, `staff_member_u`, `staff_member_d` |
| `model.ValidateFields` | `func ValidateFields(action byte, f Fielder) error` |
| `orm.DB.Tx` | `func (d *DB) Tx(fn func(tx *DB) error) error` |

## Explicit deletions (grep must be empty)

- `ops.go`: `OpListStaff`, `OpGetStaff`, `OpUpsertStaff`, `OpDeleteStaff`, raw `reg.Operation` calls, `var _ router.OperationModule`, imports `webtyp.com/router`, `webtyp.com/rpc`, `webtyp.com/view`.
  `grep -rn "reg.Operation\|router.Context\|rpc\.\|args.TenantId" --include=*.go . | grep -v _test` → empty.
- `model.go`: manual argument models for ops if any; errors covered by schema/crud.
- `module.go`: operations logic moved to `crud.go` and `queries.go`.

## Stage 1 — `model.go` and `errors.go`

- `tenant_id` and `updated_at` get `ServerOnly: true`.
- Move `Err*` constants to a new `errors.go`, Spanish doc comments with usage example each. Run `ormc`.

## Stage 2 — `crud.go` (writes) and `queries.go` (reads); `module.go` keeps core `Deps`, `Module`, `New`

- `crud.go`: implement Create, Update, Delete for `StaffMember` wrapped in `db.Tx` where appropriate, filling `tenant_id` from `m.tenantID`.
- `queries.go`: implement Read with `crud.ReadArgs` (supports ID lookup, Search, After, Limit). Always scope by `tenant_id`.

## Stage 3 — `ops.go`

- `MountOperations(reg crud.Registry)` mounts operations via generated `MountStaffMemberRead`, `MountStaffMemberCreate`, `MountStaffMemberUpdate`, `MountStaffMemberDelete`.
- `Module` implements `crud.Module` (`model.ModuleNaming` + `MountOperations(reg crud.Registry)`).

## Stage 4 — `view.go`, `seed/demo.go`, delete `web/browser.go`, `go.mod`

- `view.go`: add `//go:build browser`. Update `view.Ops` to use `model.Action`:
  `view.Ops{Module: ModelName, Allow: model.Create | model.Read | model.Update | model.Delete, Page: view.Page{First: view.All}}`
- `web/browser.go`: delete file.
- `seed/demo.go`: provide `Demo() (view.Demo, error)` returning `view.Demo{Label: ui.DefaultLabel, Timezone: "America/Santiago", ActingAs: "u1", Modules: []crud.Module{mod}, Screens: func(d view.BrowserDeps) []view.Screen { return ui.Screens(d.Caller) }}`.
- `go mod tidy`: ensure `webtyp.com/layout` is not present in `go.mod`.

## Stage 5 — tests (`tests/`)

- Update tests to use `crud.Loopback` / `crud.LoopbackActingAs("u1", mod)` and `crud.CreateArgs`, `crud.ReadArgs`, `crud.UpdateArgs`, `crud.DeleteArgs`.
- Operation names derived from `crud.Op(&staffmanager.StaffMember{}, action)` qualified with `staffmanager.ModelName + "."`.
- Add `TestDemo_Valid` verifying `seed.Demo()`.

## Stage 6 — translations and documentation

- `langc sync .`: zero untranslated, zero warnings.
- Update `README.md` and `docs/ARCHITECTURE.md`.

## Delivery check

`gotest` passes; `langc sync .` clean; every grep above empty.
