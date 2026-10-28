package tests

import (
	"testing"

	staffmanager "github.com/veltylabs/staff_manager"
)

func TestItemsDeclaredOnFields(t *testing.T) {
	rec := staffmanager.StaffMember{Id: "1", Name: "Dr. Smith", Rut: "12.345.678-9"}
	row := rec.Item()
	if row.ID != "1" || row.Label != "Dr. Smith" || row.Description != "12.345.678-9" {
		t.Errorf("Item mismatch: %+v", row)
	}
}

func TestDomainError(t *testing.T) {
	if staffmanager.ErrNotFound.Error() != "staff not found" {
		t.Errorf("expected ErrNotFound.Error() to be 'staff not found', got %q", staffmanager.ErrNotFound.Error())
	}
}
