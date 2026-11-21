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

// UserByStaff es el inverso de StaffByUser: el aviso «tu agenda cambió» se manda a la cuenta del profesional.
func TestUserByStaff(t *testing.T) {
	m := setup(t, nil)
	if err := m.RegisterFirstAdmin("user-7", "11111111-1", "Carla", "1.1.1.1"); err != nil {
		t.Fatal(err)
	}
	id, _, _, _, _ := m.StaffByUser(testTenant, "user-7")
	userID, found, err := m.UserByStaff(testTenant, id)
	if err != nil || !found || userID != "user-7" {
		t.Fatalf("UserByStaff = %q, %v, %v; want user-7, true, nil", userID, found, err)
	}
	if _, found, _ := m.UserByStaff(testTenant, "no-such-staff"); found {
		t.Error("a missing professional must not be found")
	}
}
