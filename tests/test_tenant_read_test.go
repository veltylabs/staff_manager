package tests

import (
	"github.com/veltylabs/staff_manager"
	"testing"
	"webtyp.com/model"
)

func TestTenant_ReadOtherTenantIsNotFound(t *testing.T) {
	m := setup(t, nil)
	s := staffmanager.StaffMember{Id: "s1", TenantId: "other", Rut: "111", Name: "A"}
	m.ImportStaffMember(&s)
	_, err := m.GetStaff("other", "s1")
	if model.KindOf(err) != model.ErrorNotFound {
		t.Errorf("expected NotFound, got %v", err)
	}
}
