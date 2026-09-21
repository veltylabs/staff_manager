---
PLAN: "fix: UpsertStaff writes to the database without validating the record first, unlike every sibling module"
EXECUTOR: jules
REVIEWER: none
---

> This plan is dispatched via the CodeJob workflow. See skill: agents-workflow.

# PLAN — `UpsertStaff` must validate before it writes

You are an external agent with **zero prior context** about this project. Everything you need is
in this file. Read `AGENTS.md` at the repo root first — it defines the conventions referenced
below (validate-before-write, `storage/mem` test backend, `tests/` layout) — then read this file
fully before writing code.

## 0. Prerequisite — run this first

```bash
go install webtyp.com/devflow/cmd/gotest@latest
```

All tests run with `gotest`, never `go test` directly.

## 1. The bug, proven by a failing test already in this repo

`tests/staff_test.go` now ends with `TestUpsertStaff_RejectsMissingUserID` (already committed on
this branch). Run it and confirm it is RED:

```bash
gotest
# --- FAIL: TestUpsertStaff_RejectsMissingUserID
#     staff_test.go:224: UpsertStaff with an empty UserId returned no error — ...
```

This is the acceptance criterion: **that test must go GREEN, and no other test may regress.**

## 2. Root cause

`StaffMemberModel` (`model.go`) declares:

```go
{Name: "user_id", Type: model.Text(), NotNull: true},
```

but `Module.UpsertStaff` (`module.go`) never calls `member.Validate(...)` (the generated method
from `model_orm.go`) before either `m.db.Update(...)` or `m.db.Create(&member)`. This breaks
`AGENTS.md`'s own rule:

> Every create/update path calls the generated `Validate(action)` (or `model.ValidateFields`)
> **before** `db.Create`/`db.Update` — fail-closed: data that was never validated never reaches
> the DB.

The sibling module `github.com/veltylabs/item_catalog` follows this correctly — its `mcp.go`
calls `item.Validate(action)` (and `spec.Validate(action)`, `a.Validate(action)`) immediately
before every `db.Create`/`db.Update` call. `staff_manager` is the odd one out.

**Why this went unnoticed by this module's own test suite:** every test here (per this repo's
`AGENTS.md`) builds `*orm.DB` over `webtyp.com/storage/mem`, which does **not** enforce
`NotNull` the way a real SQL backend does — an invalid row is silently accepted. In production,
the same call reaches Postgres, whose real `NOT NULL` constraint on the `user_id` column rejects
the insert — but as a **raw, untyped database error** surfacing from deep inside `db.Create`,
never as a clean domain validation error the caller (or an app's UI) can present to a person.

**Downstream symptom this caused (informational, not part of this plan):** a downstream app's
"add funcionario" screen (crudview-generated form over `StaffMemberModel`) has no field for
`user_id` — correctly, per this repo's own widget-assignment rule ("Base kinds (`model.X()`) on
ids, `tenant_id`, timestamps... — output is never rendered as an editable form"). Every save
attempt from that screen therefore submits `UserId == ""`, and today that failure is invisible
until it hits Postgres, and even then arrives as an opaque DB error with no typed shape a UI
could distinguish from any other 500.

## 3. The fix

In `module.go`, `func (m *Module) UpsertStaff(member StaffMember) (StaffMember, error)`:

Call `member.Validate(...)` right after `member.UpdatedAt = time.Now()` (once the record is
fully populated, before either branch decides update-vs-create), using the correct action per
branch:

- The "update by Id" branch and the "update by tenant+rut" branch: `model.ActionUpdate`.
- The "create new" branch (falls through both `orm.ErrNotFound` checks): `model.ActionCreate`.

Concretely, structure it as: validate with `model.ActionUpdate` before attempting either update
branch (an update failing validation must return the error immediately, exactly like the existing
`err != orm.ErrNotFound` checks already do), and validate with `model.ActionCreate` right before
`if member.Id == "" { member.Id = m.ids.NewID() }` / `m.db.Create(&member)` at the bottom of the
function, since `Id`/`TenantId` are only fully resolved by that point on the create path.

```go
// after member.UpdatedAt = time.Now()

// ... existing "update by Id" branch:
if member.Id != "" {
	var existing StaffMember
	qb := m.db.Query(&existing).Where("id").Eq(member.Id).Where("tenant_id").Eq(member.TenantId)
	err := qb.ReadOne()
	if err == nil {
		existing.Name = member.Name
		existing.Rut = member.Rut
		existing.Specialty = member.Specialty
		existing.Role = member.Role
		existing.IsActive = member.IsActive
		existing.UserId = member.UserId
		existing.UpdatedAt = member.UpdatedAt
		if err := existing.Validate(model.ActionUpdate); err != nil {
			return StaffMember{}, err
		}
		if err := m.db.Update(&existing, orm.Eq("id", existing.Id), orm.Eq("tenant_id", existing.TenantId)); err != nil {
			return StaffMember{}, err
		}
		return existing, nil
	}
	if err != orm.ErrNotFound {
		return StaffMember{}, err
	}
}
// ... existing "update by tenant+rut" branch: same shape — validate the
// mutated `existing` with model.ActionUpdate right before m.db.Update.

// ... existing "create new" branch, right before `if member.Id == "" { ... }`:
if member.Id == "" {
	member.Id = m.ids.NewID()
}
if err := member.Validate(model.ActionCreate); err != nil {
	return StaffMember{}, err
}
if err := m.db.Create(&member); err != nil {
	return StaffMember{}, err
}
return member, nil
```

Import `webtyp.com/model` in `module.go` if it is not already imported under that name (it is
already imported for `model.IDGenerator`, so no new import line should be needed — just use the
existing alias).

### What NOT to do

- **Do not add a manual `if member.UserId == "" { return ... }` check.** `AGENTS.md` is explicit:
  "Manual `if x == ""` checks may exist as defense in depth but never replace the declared
  constraint." The fix is to call the existing generated `Validate`, which already enforces every
  `NotNull` field declared in `StaffMemberModel` — not to hand-roll a check for this one field.
- **Do not change `StaffMemberModel`** (e.g. do not drop `NotNull: true` from `user_id`, and do
  not add an `input.*` widget to it — per `AGENTS.md`'s widget-assignment rule, `user_id` is not
  user-editable and must stay a plain `model.Text()`).
- **Do not design or implement a "create a linked auth user for a new staff member" flow.** That
  is a real, separate product decision (who provisions `user_id` for a brand-new funcionario, and
  how) that has not been made yet and is explicitly **out of scope** for this plan. This plan only
  makes the existing, already-declared constraint fail loudly and cleanly instead of silently or
  as a raw DB error — it does not change what values are valid.
- **Do not touch `storage/mem`** to make it enforce `NotNull`. That is a `webtyp.com/storage`
  concern (if it is even desired there — the divergence is worth flagging upstream separately),
  not this module's.

## 4. Verification

```bash
gotest
# vet ✅, race ✅, tests ✅ — TestUpsertStaff_RejectsMissingUserID now PASSES,
# and every other existing test in tests/staff_test.go, tests/trust_test.go still passes
# (TestStaffExists_True and the two update-path tests all supply a non-empty UserId already,
# so they are unaffected).
```

## Stages

| # | Stage | File(s) | Acceptance |
|---|---|---|---|
| 1 | Validate before the "update by Id" write | `module.go` | `existing.Validate(model.ActionUpdate)` called, error returned immediately on failure |
| 2 | Validate before the "update by tenant+rut" write | `module.go` | Same, for that branch's `existing` |
| 3 | Validate before the "create" write | `module.go` | `member.Validate(model.ActionCreate)` called after `Id` is assigned, before `db.Create` |
| 4 | Verify | — | `gotest` green, `TestUpsertStaff_RejectsMissingUserID` passes, no other test changed |
