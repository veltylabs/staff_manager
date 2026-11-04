---
PLAN: "feat: StaffByUser for clinical_encounter; required contracts move to interfaces.go"
EXECUTOR: jules
REVIEWER: none
STATUS: running
SESSION: 2264374540799879527
---

> Read AGENTS.md at the repo root before any change.
> **Finish by yourself:** when the delivery check passes, commit, push and open the pull request. Never ask whether to
> proceed or finalize — nobody answers, and a session that waits ends without your work.

# Plan — staff_manager answers "which professional is this user" and keeps its contracts in interfaces.go

## Tools

```bash
go install webtyp.com/devflow/cmd/gotest@latest
go get webtyp.com/time@v0.6.4 webtyp.com/appkit@v0.2.8
```

## Why

`clinical_encounter` (v0.2.53) no longer imports this module: its `module/server.go` finds, with
`appkit.Find`, a value with this method (from its `interfaces.go`):

```go
type DoctorReader interface {
	StaffByUser(tenantID, userID string) (staffID, name, specialty string, found bool, err error)
}
```

Nobody provides it, so an app with `clinical_encounter` does not start. Ecosystem rule (module-layout): the contracts a
module requires from others live only in `interfaces.go`; `module.go` keeps `Deps`, `Module` and `New`.

## Ecosystem search

`StaffExists`, `StaffCard`, `GetStaff` (`module.go`) exist → reuse their style. `StaffMember` has `user_id`
(`model.go:13`, text, not null). New: `StaffByUser` (no equivalent).

## Design gate

1. Prior art: same narrow-reader shape as `StaffCard`, satisfied structurally (Go interfaces; no adapter).
2. Name: `StaffByUser` is the name the consumer already declared.
3. Ledger: +1 method; 0 new ways to do the same thing.
4. Belongs here: staff_manager owns the staff ↔ user link.
5. Deletes: nothing; moves `DeviceReader` and `DeviceRegistrar` from `module.go` to `interfaces.go`.

## Stage 1 — interfaces.go

Move `DeviceReader` and `DeviceRegistrar` (with their comments) from `module.go` (lines ~11–21) to a new
`interfaces.go` (package `staffmanager`), unchanged.

## Stage 2 — StaffByUser (`module.go`, next to `StaffCard`)

```go
// StaffByUser devuelve el profesional cuya cuenta de usuario es userID. Satisface
// DoctorReader de clinical_encounter estructuralmente. Sin fila ⇒ found=false y err=nil;
// solo una falla de la base devuelve error.
// Ejemplo: staffID, name, spec, ok, err := m.StaffByUser("t1", "u1")
func (m *Module) StaffByUser(tenantID, userID string) (staffID, name, specialty string, found bool, err error) {
	if tenantID == "" || userID == "" {
		return "", "", "", false, nil
	}
	qb := m.db.Query(&StaffMember{}).Where(StaffMember_.TenantId).Eq(tenantID).Where(StaffMember_.UserId).Eq(userID)
	rows, err := ReadAllStaffMember(qb)
	if err != nil && !orm.IsNotFound(err) {
		return "", "", "", false, err
	}
	if len(rows) == 0 {
		return "", "", "", false, nil
	}
	s := rows[0]
	return s.Id, s.Name, s.Specialty, true, nil
}
```

If `StaffMember_.UserId`/`StaffMember_.TenantId` are named differently in `model_orm.go`, use the generated names.

## Stage 3 — tests (`tests/`, sqlite like the existing ones)

| Test | Setup → call | Expected |
|---|---|---|
| `TestStaffByUserFound` | staff "s1" with user_id "u1", name "Ana", specialty "Pediatría" → `StaffByUser("t1","u1")` | `"s1","Ana","Pediatría",true,nil` |
| `TestStaffByUserMissing` | no staff → `StaffByUser("t1","u9")` | `found == false`, `err == nil` |
| `TestStaffByUserOtherTenant` | staff of "t2" with user "u1" → `StaffByUser("t1","u1")` | `found == false` |

## Stage 4 — docs (Spanish)

`README.md` (Step-by-Step) and `docs/ARCHITECTURE.md`: list `StaffByUser` among the readers it offers to other
modules, and say the required contracts are in `interfaces.go`.

## Delivery check

1. `gotest` passes.
2. `grep -n "interface {" module.go` → empty.

Commit after every stage. Never edit `model_orm.go` by hand; a build error there means the hand-written code is wrong.
Never run `git reset`, `git restore`, `git checkout -- <file>` or `git stash` on your own work. Stuck 30 minutes: open
the PR anyway listing what fails. A PR that does not compile, or whose tests could not run, is a failed delivery.
Environment: `GOPRIVATE=webtyp.com/*,github.com/veltylabs/*`, `GOFLAGS=-mod=mod`, `go mod tidy` before `ormc`.
If a signature in the published tag differs from the APIs table, stop and name the symbol in the PR; never adapt it by guessing.

## Stages

| Stage | Files | Done when |
|---|---|---|
| 1 | `interfaces.go`, `module.go` | compiles |
| 2 | `module.go` | compiles |
| 3 | `tests/staff_by_user_test.go` | tests pass |
| 4 | `README.md`, `docs/ARCHITECTURE.md` | done |
