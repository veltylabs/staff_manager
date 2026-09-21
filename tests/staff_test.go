package tests

import (
	"testing"

	staffmanager "github.com/veltylabs/staff_manager"
)

// El puerto que appointment_booking declara en su propio lado. Redeclarado localmente a
// propósito: este repositorio no debe depender de un módulo de programación para probar que
// satisface una interfaz estructural.
type staffReader interface {
	StaffExists(tenantId, staffId string) (bool, error)
}

var _ staffReader = (*staffmanager.Module)(nil)

func TestStaffExists_True(t *testing.T) {
	m := setup(t, nil)
	member, err := m.UpsertStaff(staffmanager.StaffMember{
		UserId: "user-1",
		Rut:    "12345678-5",
		Name:   "Dr. Smith",
	})
	if err != nil {
		t.Fatalf("UpsertStaff: %v", err)
	}

	exists, err := m.StaffExists(testTenant, member.Id)
	if err != nil {
		t.Fatalf("StaffExists error: %v", err)
	}
	if !exists {
		t.Errorf("StaffExists(%q, %q) = false, want true", testTenant, member.Id)
	}
}

func TestStaffExists_UnknownID(t *testing.T) {
	m := setup(t, nil)
	exists, err := m.StaffExists(testTenant, "non-existent-id")
	if err != nil {
		t.Fatalf("StaffExists error: %v", err)
	}
	if exists {
		t.Errorf("StaffExists = true, want false for unknown ID")
	}
}

func TestStaffExists_WrongTenant(t *testing.T) {
	m := setup(t, nil)
	member, err := m.UpsertStaff(staffmanager.StaffMember{
		UserId: "user-1",
		Rut:    "12345678-5",
		Name:   "Dr. Smith",
	})
	if err != nil {
		t.Fatalf("UpsertStaff: %v", err)
	}

	exists, err := m.StaffExists("other-tenant", member.Id)
	if err != nil {
		t.Fatalf("StaffExists error: %v", err)
	}
	if exists {
		t.Errorf("StaffExists = true, want false for wrong tenant")
	}
}

func TestStaffExists_EmptyArgs(t *testing.T) {
	m := setup(t, nil)

	cases := []struct {
		tenantID string
		staffID  string
	}{
		{"", ""},
		{testTenant, ""},
		{"", "staff-1"},
	}

	for _, c := range cases {
		exists, err := m.StaffExists(c.tenantID, c.staffID)
		if err != nil {
			t.Fatalf("StaffExists(%q, %q) returned error: %v", c.tenantID, c.staffID, err)
		}
		if exists {
			t.Errorf("StaffExists(%q, %q) = true, want false", c.tenantID, c.staffID)
		}
	}
}

func TestUpsertStaff_PersistsSpecialtyAndRole(t *testing.T) {
	m := setup(t, nil)
	created, err := m.UpsertStaff(staffmanager.StaffMember{
		UserId:    "user-1",
		Rut:       "12345678-5",
		Name:      "Dr. Smith",
		Specialty: "Radiology",
		Role:      "Radiologist",
	})
	if err != nil {
		t.Fatalf("UpsertStaff: %v", err)
	}

	got, err := m.GetStaff(testTenant, created.Id)
	if err != nil {
		t.Fatalf("GetStaff: %v", err)
	}
	if got.Specialty != "Radiology" {
		t.Errorf("Specialty = %q, want %q", got.Specialty, "Radiology")
	}
	if got.Role != "Radiologist" {
		t.Errorf("Role = %q, want %q", got.Role, "Radiologist")
	}
}

func TestUpsertStaff_UpdateByIDKeepsSpecialty(t *testing.T) {
	m := setup(t, nil)
	created, err := m.UpsertStaff(staffmanager.StaffMember{
		UserId:    "user-1",
		Rut:       "12345678-5",
		Name:      "Dr. Smith",
		Specialty: "Radiology",
		Role:      "Radiologist",
	})
	if err != nil {
		t.Fatalf("UpsertStaff: %v", err)
	}

	// Ruta de actualización A: por ID
	updated, err := m.UpsertStaff(staffmanager.StaffMember{
		Id:        created.Id,
		UserId:    "user-1",
		Rut:       "12345678-5",
		Name:      "Dr. Smith",
		Specialty: "Cardiology",
		Role:      "Cardiologist",
	})
	if err != nil {
		t.Fatalf("UpsertStaff (update by ID): %v", err)
	}

	got, err := m.GetStaff(testTenant, updated.Id)
	if err != nil {
		t.Fatalf("GetStaff: %v", err)
	}
	if got.Specialty != "Cardiology" {
		t.Errorf("Specialty = %q, want %q", got.Specialty, "Cardiology")
	}
	if got.Role != "Cardiologist" {
		t.Errorf("Role = %q, want %q", got.Role, "Cardiologist")
	}
}

func TestUpsertStaff_UpdateByRutKeepsSpecialty(t *testing.T) {
	m := setup(t, nil)
	_, err := m.UpsertStaff(staffmanager.StaffMember{
		UserId:    "user-1",
		Rut:       "12345678-5",
		Name:      "Dr. Smith",
		Specialty: "Radiology",
		Role:      "Radiologist",
	})
	if err != nil {
		t.Fatalf("UpsertStaff: %v", err)
	}

	// Ruta de actualización B: sin ID, coincidiendo tenant+RUT
	updated, err := m.UpsertStaff(staffmanager.StaffMember{
		UserId:    "user-1",
		Rut:       "12345678-5",
		Name:      "Dr. Smith",
		Specialty: "Neurology",
		Role:      "Neurologist",
	})
	if err != nil {
		t.Fatalf("UpsertStaff (update by RUT): %v", err)
	}

	got, err := m.GetStaff(testTenant, updated.Id)
	if err != nil {
		t.Fatalf("GetStaff: %v", err)
	}
	if got.Specialty != "Neurology" {
		t.Errorf("Specialty = %q, want %q", got.Specialty, "Neurology")
	}
	if got.Role != "Neurologist" {
		t.Errorf("Role = %q, want %q", got.Role, "Neurologist")
	}
}

// TestUpsertStaff_RejectsMissingUserID reproduces the bug found while
// manually testing the "Funcionarios" screen of a downstream app
// (mjosefa-cms): the crudview-generated form for StaffMember never collects
// or generates user_id (it has no input.* widget — see AGENTS.md's "Widgets
// are assigned by ROLE" rule, which correctly keeps it a plain model.Text()
// since it is not user-editable), so every "add funcionario" attempt through
// that GUI submits a StaffMember with UserId == "". StaffMemberModel
// declares `{Name: "user_id", ..., NotNull: true}`, but UpsertStaff (module.go)
// never calls member.Validate(...) before db.Create/db.Update — unlike every
// sibling module (see github.com/veltylabs/item_catalog's mcp.go, which calls
// item.Validate(action) before every Create/Update, per AGENTS.md's "Every
// create/update path calls the generated Validate(action) ... BEFORE
// db.Create/db.Update — fail-closed").
//
// Because storage/mem (this test's own backend) does not enforce NotNull the
// way Postgres does, this defect was invisible to this module's own test
// suite: UpsertStaff silently "succeeds" here with an empty UserId, while in
// production the same call reaches Postgres's real NOT NULL constraint and
// fails with a raw, untyped database error deep inside db.Create — never a
// clean domain validation error the caller (or a UI) could show. This test
// must fail today and pass once UpsertStaff validates before writing.
func TestUpsertStaff_RejectsMissingUserID(t *testing.T) {
	m := setup(t, nil)

	_, err := m.UpsertStaff(staffmanager.StaffMember{
		Rut:      "12345678-5",
		Name:     "Ana Torres",
		IsActive: true,
		// UserId deliberately left empty — this is exactly what the GUI's
		// auto-save form sends today, since it has no field for it.
	})
	if err == nil {
		t.Fatal("UpsertStaff with an empty UserId returned no error — StaffMemberModel declares " +
			"user_id NotNull:true, but nothing validates that before the row reaches the database. " +
			"In production (Postgres, which DOES enforce NOT NULL, unlike this test's storage/mem " +
			"backend) this same call fails with a raw DB error instead of a clean validation error")
	}
}
