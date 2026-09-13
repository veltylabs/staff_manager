package tests

import (
	"testing"

	staffmanager "github.com/veltylabs/staff_manager"
)

// The port appointment_booking declares on its own side. Redeclared locally on
// purpose: this repository must not depend on a scheduling module to prove it
// satisfies a structural interface.
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

	// Update path A: by ID
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

	// Update path B: without ID, matching tenant+RUT
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
