package tests

import (
	"testing"

	staffmanager "github.com/veltylabs/staff_manager"
)

// TestIsTrustedIP is the consumer-shaped proof for the trust decision: it
// must answer strictly through the injected DeviceReader, never by reading
// device_manager's own table directly (that reimplementation — deviceProbe —
// was removed; see docs/PLAN.md in veltylabs/mjosefa-cms for why).
func TestIsTrustedIP(t *testing.T) {
	devices := &fakeDevices{byIP: map[string]string{
		"192.168.1.10": "device-1",
		"192.168.1.20": "device-2",
	}}
	m := setup(t, devices)

	member, err := m.UpsertStaff(staffmanager.StaffMember{UserId: "user-1", Rut: "12345678-5", Name: "Ana", IsActive: true})
	if err != nil {
		t.Fatalf("UpsertStaff: %v", err)
	}
	if err := m.AssignDevice(member.Id, "device-1"); err != nil {
		t.Fatalf("AssignDevice: %v", err)
	}

	otherMember, err := m.UpsertStaff(staffmanager.StaffMember{UserId: "user-2", Rut: "11111111-1", Name: "Beto", IsActive: true})
	if err != nil {
		t.Fatalf("UpsertStaff (other): %v", err)
	}
	if err := m.AssignDevice(otherMember.Id, "device-2"); err != nil {
		t.Fatalf("AssignDevice (other): %v", err)
	}

	inactiveMember, err := m.UpsertStaff(staffmanager.StaffMember{UserId: "user-3", Rut: "22222222-2", Name: "Cami", IsActive: false})
	if err != nil {
		t.Fatalf("UpsertStaff (inactive): %v", err)
	}
	if err := m.AssignDevice(inactiveMember.Id, "device-1"); err != nil {
		t.Fatalf("AssignDevice (inactive): %v", err)
	}

	noDeviceMember, err := m.UpsertStaff(staffmanager.StaffMember{UserId: "user-4", Rut: "33333333-3", Name: "Deni", IsActive: true})
	if err != nil {
		t.Fatalf("UpsertStaff (no device): %v", err)
	}

	cases := []struct {
		name   string
		userID string
		ip     string
		want   bool
	}{
		{"assigned device + active member", "user-1", "192.168.1.10", true},
		{"device belongs to another member", "user-1", "192.168.1.20", false},
		{"unknown ip", "user-1", "10.0.0.1", false},
		{"deactivated member", "user-3", "192.168.1.10", false},
		{"member with no devices", noDeviceMember.UserId, "192.168.1.10", false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := m.IsTrustedIP(c.userID, c.ip); got != c.want {
				t.Errorf("IsTrustedIP(%q, %q) = %v, want %v", c.userID, c.ip, got, c.want)
			}
		})
	}
}
