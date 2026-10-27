package tests

import (
	"testing"

	"webtyp.com/model"

	staffmanager "github.com/veltylabs/staff_manager"
)

func TestRowsDeclaredOnFields(t *testing.T) {
	rec := staffmanager.StaffMember{Id: "1", Name: "Dr. Smith", Rut: "12.345.678-9"}
	row := model.RowOf(&rec)
	if row.ID != "1" || row.Title != "Dr. Smith" || row.Detail != "12.345.678-9" {
		t.Errorf("RowOf mismatch: %+v", row)
	}
}

func TestDomainError(t *testing.T) {
	if staffmanager.ErrNotFound.Error() != "staff not found" {
		t.Errorf("expected ErrNotFound.Error() to be 'staff not found', got %q", staffmanager.ErrNotFound.Error())
	}
}
