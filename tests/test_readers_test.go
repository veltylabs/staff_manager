package tests

import (
	"testing"
)

func TestReaders_RegisterFirstAdmin(t *testing.T) {
	m := setup(t, nil)
	err := m.RegisterFirstAdmin("user-2", "11111111-1", "Bob", "1.1.1.1")
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}
	id, _, _, found, _ := m.StaffByUser(testTenant, "user-2")
	if !found {
		t.Errorf("expected found")
	}
	staff, _ := m.GetStaff(testTenant, id)
	if staff.Name != "Bob" {
		t.Errorf("expected Bob, got %v", staff.Name)
	}
}
