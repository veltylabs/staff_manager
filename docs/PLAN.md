---
PLAN: "feat: specialty and role on a staff member, plus the StaffExists port"
EXECUTOR: jules
REVIEWER: none
STATUS: review
SESSION: 17733230379642182941
PR: https://github.com/veltylabs/staff_manager/pull/1
---

> This plan is dispatched via the CodeJob workflow. See skill: **agents-workflow**.

# Plan — `specialty`, `role`, and `StaffExists`

You are an agent with **no prior context** and you have **only this repository**
(`github.com/veltylabs/staff_manager`). Everything you need is inline.

## 1. The problem

Two gaps, both discovered while wiring this module next to
`github.com/veltylabs/appointment_booking` in a clinic application.

### 1.1 Nothing can ask whether a staff member exists

`appointment_booking` refuses to create a reservation for a professional it
cannot verify. It declares the question as a narrow port it does not import:

```go
// appointment_booking/service.go
// StaffReader verifica que un miembro del staff existe y pertenece al tenant.
type StaffReader interface {
	StaffExists(tenantId, staffId string) (bool, error)
}
```

and its own README documents the expected wiring as
`Staff: staffmodule.New(db)  // implements StaffReader` — that is, the staff
module itself satisfies the port structurally, with no adapter.

`*Module` does not have that method. Every application wanting to book an
appointment must therefore write the same six-line adapter over `GetStaff`,
comparing `ErrNotFound` by hand. That adapter is application code doing a
module's job, and it will be written once per application.

### 1.2 A staff member has no specialty and no role

`StaffMemberModel` today carries `id`, `tenant_id`, `user_id`, `rut`, `name`,
`is_active`, `updated_at`. There is nowhere to record that someone is a
radiologist.

The consequence is already visible in the ecosystem: the sibling module
`github.com/veltylabs/clinical_encounter` writes a
`doctor_specialty_snapshot` column on every medical record, and there is no
source for that value — it is typed by hand at the point of use, per record,
with nothing keeping two spellings of the same specialty together. A booking
screen that wants to offer "who does ultrasound" has nothing to group by.

## 2. Design gate

### Prior art

- **Django** `Model.objects.filter(...).exists()` and **Rails**
  `Model.exists?(id)` both expose existence as a first-class boolean query,
  separate from fetching the row, precisely because callers that only need the
  answer should not pay for or handle a full record.
- **Ent** (Go) generates `Query().Exist(ctx)` per entity for the same reason.
- All three return only a boolean plus an error. We match that shape exactly.
  We differ in one way: ours is tenant-scoped (`tenantId` is the first
  argument), because in this ecosystem no read is ever global — every existing
  method on this module already takes the tenant first.

For the two new fields there is no interesting prior art to differ from: they
are plain columns on the module's own record.

### Novice-name test

- "Does this staff exist in this tenant?" → `StaffExists(tenantId, staffId)`.
  The port `appointment_booking` declares already uses this exact name and
  signature; choosing a different one would mean every consumer writes an
  adapter, which is the defect being fixed.
- "The staff member's specialty" → `member.Specialty`.
- "The staff member's role" → `member.Role`.

### Complexity ledger

| | change |
|---|---|
| Concepts | `+0` — existence and specialty are both already concepts in this domain; neither introduces a new one |
| Files | `+1` (`tests/staff_test.go`); `model.go`, `module.go`, `migrate/migrate.go` are edited |
| Call-site lines | `−6 per consuming application` (the hand-written `StaffReader` adapter disappears), `+1` here |
| Ways to do it | `−1` — "is this professional real?" had zero supported answers and one copied one; now it has one |
| Net | negative |

### Where it belongs

Here. `StaffExists` is a question about this module's own table, and the two
fields are attributes of this module's own record. Putting the specialty
anywhere else — a lookup table in the catalog, a column on the clinical
record — is what produced the unsourced `doctor_specialty_snapshot` in the
first place.

### What it deletes

Nothing in this repository; this is additive. What it makes deletable is the
per-application `StaffReader` adapter, in those applications.

## 3. Decisions already taken — do not revisit

1. **`specialty` is free text (`input.Text()`), not a slug and not an enum.**
   A slug column would need a slug→label table, and the only existing instance
   of that in this ecosystem is a hardcoded Spanish `switch` in a demo — an
   application hardcoding data the database could supply. A consumer that wants
   to group professionals by specialty groups by the distinct values it reads
   back from `list_staff`. This module renders no human language of its own.
2. **`role` is free text too**, and it is **not** a permission. Authorization in
   this ecosystem is `webtyp.com/rbac`, keyed by user, resource and action.
   `role` here is a descriptive job title shown to a person ("Médico",
   "TENS", "Administrativo"). Do not wire it to anything that grants access.
3. **Both fields are optional** (no `NotNull`). Existing rows have neither, and
   a required column would break every current record.
4. **`Item()` is not touched.** `StaffMember.Item()` currently returns
   `Description: m.Rut`, and an existing staff screen renders that. Changing it
   to show the specialty is a UI decision for the consuming application, not
   this plan.

## 4. Stages

### Stage 1 — the two fields, in `model.go`

Add both to `StaffMemberModel.Fields`, **immediately after `name`** and before
`is_active`:

```go
{Name: "specialty", Type: input.Text(), Permitted: model.Permitted{Maximum: 120}},
{Name: "role", Type: input.Text(), Permitted: model.Permitted{Maximum: 60}},
```

Write a comment above them recording decision 3.1 and 3.2 from this plan —
that the values are free text on purpose, and that `role` is descriptive and
never an authorization input. That reasoning must outlive this plan file.

Regenerate `model_orm.go` with `ormc`. **Never hand-edit `model_orm.go`.** The
generated struct gains `Specialty` and `Role` (pure casing of the column
names).

### Stage 2 — carry the fields through `UpsertStaff`, in `module.go`

`UpsertStaff` has **two** update paths, and both copy fields field-by-field
onto an `existing` record. A new column added to only one of them silently
disappears on the other path.

**Path A — update by id** (`if member.Id != ""`, after the `ReadOne` succeeds):

```go
existing.Name = member.Name
existing.Rut = member.Rut
existing.Specialty = member.Specialty   // ADD
existing.Role = member.Role             // ADD
existing.IsActive = member.IsActive
existing.UserId = member.UserId
existing.UpdatedAt = member.UpdatedAt
```

**Path B — update by tenant+rut** (the block after
`Where("tenant_id").Eq(...).Where("rut").Eq(...)` succeeds):

```go
existing.Name = member.Name
existing.Specialty = member.Specialty   // ADD
existing.Role = member.Role             // ADD
existing.IsActive = member.IsActive
existing.UpdatedAt = member.UpdatedAt
```

The create path needs no change — it writes `member` whole.

**Acceptance:** `grep -n "existing.Specialty" module.go` → exactly 2 matches.
Same for `existing.Role`.

### Stage 3 — `StaffExists`, in `module.go`

Add immediately after `GetStaff`, since it is the same query narrowed to a
boolean:

```go
// StaffExists reports whether a staff member with this id belongs to this
// tenant. It satisfies the narrow StaffReader port that a scheduling module
// declares on its own side (StaffExists(tenantId, staffId) (bool, error)) —
// structurally, with no adapter and no import in either direction.
//
// A missing row is (false, nil), NOT an error: "this id is not one of ours" is
// the answer the caller asked for. Only a real storage failure returns a
// non-nil error, so a caller can never mistake a dead database for a clean
// "no".
func (m *Module) StaffExists(tenantID, staffID string) (bool, error) {
	if tenantID == "" || staffID == "" {
		return false, nil
	}
	_, err := m.GetStaff(tenantID, staffID)
	if err != nil {
		if err == ErrNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
```

**Anti-footgun.** Do not "optimise" this by inlining the query and returning
`err == nil`. `GetStaff` is the single place that maps `orm.ErrNotFound` to the
domain sentinel; a second copy of that mapping is how the two drift. And do not
collapse the error to a bare `bool` — a caller that cannot tell "absent" from
"database unreachable" will book an appointment against a staff member it never
verified.

### Stage 4 — `migrate/` must ADD COLUMNS, not just create tables

This is the stage most likely to be got wrong, because the current code looks
finished.

`migrate/migrate.go` today calls `CreateTable`, which compiles to
`CREATE TABLE IF NOT EXISTS`. Against a database where `staff_member` already
exists — which is every deployment that has ever run — it is a **no-op**, so
the two new columns would never appear and every write would fail on an unknown
column.

`webtyp.com/ddl` already has the right operation: `(*ddl.DB).Sync(models...)`
emits `CreateTable` and then, for each column in the model that the table does
not have, an additive `OpAddColumn` — inside a transaction where the backend
supports one. Replace the body:

```go
func Migrate(conn ddl.Execer, ddlCompiler ddl.Compiler) error {
	// Sync, not CreateTable: CreateTable compiles to CREATE TABLE IF NOT
	// EXISTS and is a no-op against a table that already exists, so a column
	// added to a model would never reach a deployed database. Sync creates the
	// table when it is absent and adds the missing columns when it is not.
	return ddl.New(conn, ddlCompiler).Sync(
		&staffmanager.StaffMember{},
		&staffmanager.StaffDevice{},
	)
}
```

`Sync` is additive only — it never drops or narrows a column — so it is safe to
run repeatedly and safe to run against production.

**Anti-footgun.** Sibling modules (`device_manager`, `item_catalog`,
`clinical_encounter`) still call `CreateTable` in their own `migrate/`
packages. You have **only this repository**; do not attempt to change them, and
do not treat their code as the pattern to copy here. The difference is
deliberate: this is the module adding columns to an existing table.

### Stage 5 — tests, in `tests/`

`tests/` currently holds `setup_test.go` and `trust_test.go`. Read
`setup_test.go` first and reuse its module construction and its fakes — do not
build a second harness.

New file `tests/staff_test.go`:

| Test | Asserts |
|---|---|
| `TestStaffExists_True` | an upserted member returns `(true, nil)` |
| `TestStaffExists_UnknownID` | `(false, nil)` — **not** an error |
| `TestStaffExists_WrongTenant` | a member of tenant A is `(false, nil)` for tenant B |
| `TestStaffExists_EmptyArgs` | `("", "")` and `("t1", "")` both `(false, nil)` |
| `TestUpsertStaff_PersistsSpecialtyAndRole` | create, then `GetStaff`, both fields round-trip |
| `TestUpsertStaff_UpdateByIDKeepsSpecialty` | update path A changes the specialty and the read-back shows the new value |
| `TestUpsertStaff_UpdateByRutKeepsSpecialty` | update path B — upsert **without** an `Id`, matching on tenant+rut — changes the specialty and the read-back shows the new value |

The last two are what catch a field carried through only one of the two update
paths. Both must exist.

Add one compile-time assertion, in `tests/staff_test.go`, that pins the port
shape so a future signature change fails the build here rather than in a
consumer:

```go
// The port appointment_booking declares on its own side. Redeclared locally on
// purpose: this repository must not depend on a scheduling module to prove it
// satisfies a structural interface.
type staffReader interface {
	StaffExists(tenantId, staffId string) (bool, error)
}

var _ staffReader = (*staffmanager.Module)(nil)
```

Run `gotest ./...` — everything green, including `trust_test.go` untouched.

### Stage 6 — documentation

`README.md` is four lines and an Ops table. Extend it:

- State what a staff member is: the tenant-scoped registry linking a user (by
  RUT) to devices, **and** carrying their specialty and role.
- Add a **"Ports this module satisfies"** section naming `IsTrustedIP`
  (`auth.TrustedIPStore`) and the new `StaffExists`, each with its one-line
  question and its signature.
- Add a **Schema** section listing `StaffMember`'s columns, marking `specialty`
  and `role` as free text and stating that `role` is descriptive, never an
  authorization input.
- Add a note that `migrate.Migrate` uses `Sync` and is additive.

Do **not** link any permanent document to `docs/PLAN.md` — it is deleted when
this lands.

## 5. Stages table

| # | Stage | Files | Acceptance |
|---|---|---|---|
| 1 | Fields | `model.go`, `model_orm.go` (ormc) | `Specialty` and `Role` on the generated struct |
| 2 | Upsert paths | `module.go` | `grep -c "existing.Specialty" module.go` → 2 |
| 3 | `StaffExists` | `module.go` | absent row → `(false, nil)` |
| 4 | Additive migration | `migrate/migrate.go` | `grep -n "CreateTable" migrate/migrate.go` → empty |
| 5 | Tests | `tests/staff_test.go` (new) | 7 cases + the port assertion |
| 6 | Docs | `README.md` | ports and schema documented |

## 6. Acceptance criteria

- `grep -n "CreateTable" migrate/migrate.go` → **empty**; `Sync` is the only
  DDL call.
- `grep -c "existing.Specialty" module.go` → `2`; same for `existing.Role`.
- `var _ staffReader = (*staffmanager.Module)(nil)` compiles in
  `tests/staff_test.go`.
- `StaffExists` returns `(false, nil)` — never an error — for an unknown id, an
  id from another tenant, and empty arguments.
- `model_orm.go` changed only through `ormc` regeneration.
- `StaffMember.Item()` is unchanged.
- No hardcoded specialty list, slug table, or label `switch` anywhere in the
  repository: `grep -rniE "radiolog|traumatolog|kinesiolog" .` → **empty**.
- `gotest ./...` green, `tests/trust_test.go` unmodified.
