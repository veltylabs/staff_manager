package tests

import (
	"testing"

	staffmanager "github.com/veltylabs/staff_manager"
)

// importStaff carga una persona del equipo y responde su fila ya guardada.
func importStaff(t *testing.T, m *staffmanager.Module, row staffmanager.StaffMember) staffmanager.StaffMember {
	t.Helper()
	if err := m.ImportStaffMember(&row); err != nil {
		t.Fatalf("ImportStaffMember(%q): %v", row.Name, err)
	}
	return row
}

// TestIsTrustedIP es la prueba con forma de consumidor para la decisión de confianza: debe responder
// estrictamente a través del DeviceReader inyectado, nunca leyendo la tabla de device_manager.
func TestTrust_IsTrustedIPDecisions(t *testing.T) {
	devices := &fakeDevices{byIP: map[string]string{"192.168.1.10": "device-1", "192.168.1.20": "device-2"}}
	m := setup(t, devices)

	member := importStaff(t, m, staffmanager.StaffMember{UserId: "user-1", Rut: "12345678-5", Name: "Ana", IsActive: true})
	if err := m.AssignDevice(member.Id, "device-1"); err != nil {
		t.Fatalf("AssignDevice: %v", err)
	}
	other := importStaff(t, m, staffmanager.StaffMember{UserId: "user-2", Rut: "11111111-1", Name: "Beto", IsActive: true})
	if err := m.AssignDevice(other.Id, "device-2"); err != nil {
		t.Fatalf("AssignDevice (other): %v", err)
	}
	inactive := importStaff(t, m, staffmanager.StaffMember{UserId: "user-3", Rut: "22222222-2", Name: "Cami", IsActive: false})
	if err := m.AssignDevice(inactive.Id, "device-1"); err != nil {
		t.Fatalf("AssignDevice (inactive): %v", err)
	}
	noDevice := importStaff(t, m, staffmanager.StaffMember{UserId: "user-4", Rut: "33333333-3", Name: "Deni", IsActive: true})

	cases := []struct {
		name   string
		userID string
		ip     string
		want   bool
	}{
		{"assigned device and active member", "user-1", "192.168.1.10", true},
		{"device belongs to another member", "user-1", "192.168.1.20", false},
		{"unknown ip", "user-1", "10.0.0.1", false},
		{"deactivated member", "user-3", "192.168.1.10", false},
		{"member with no devices", noDevice.UserId, "192.168.1.10", false},
	}
	for _, c := range cases {
		if got := m.IsTrustedIP(c.userID, c.ip); got != c.want {
			t.Errorf("%s: IsTrustedIP(%q, %q) = %v, want %v", c.name, c.userID, c.ip, got, c.want)
		}
	}
}

func TestReaders_StaffByUser(t *testing.T) {
	m := setup(t, nil)
	member := importStaff(t, m, staffmanager.StaffMember{UserId: "u1", Name: "Ana", Specialty: "Pediatría", Rut: "1-9", IsActive: true})

	staffID, name, specialty, found, err := m.StaffByUser(testTenant, "u1")
	if err != nil || !found || staffID != member.Id || name != "Ana" || specialty != "Pediatría" {
		t.Fatalf("StaffByUser(u1) = %q %q %q found=%v err=%v", staffID, name, specialty, found, err)
	}
	if _, _, _, found, err := m.StaffByUser(testTenant, "nobody"); err != nil || found {
		t.Fatalf("an unknown user must answer found=false and no error, got found=%v err=%v", found, err)
	}
}

func TestReaders_StaffCardAndExists(t *testing.T) {
	m := setup(t, nil)
	member := importStaff(t, m, staffmanager.StaffMember{UserId: "u1", Name: "Ana", Specialty: "Pediatría", Rut: "1-9", IsActive: true})

	name, specialty, active, err := m.StaffCard(testTenant, member.Id)
	if err != nil || name != "Ana" || specialty != "Pediatría" || !active {
		t.Fatalf("StaffCard = %q %q active=%v err=%v", name, specialty, active, err)
	}
	if _, _, _, err := m.StaffCard("other-tenant", member.Id); err != staffmanager.ErrNotFound {
		t.Fatalf("a card asked for another tenant must be ErrNotFound, got %v", err)
	}
	if ok, err := m.StaffExists(testTenant, member.Id); err != nil || !ok {
		t.Fatalf("StaffExists = %v err=%v", ok, err)
	}
	if ok, err := m.StaffExists(testTenant, "missing"); err != nil || ok {
		t.Fatalf("a missing person must answer false and no error, got %v err=%v", ok, err)
	}
}

func TestReaders_ListStaffSearchesByName(t *testing.T) {
	m := setup(t, nil)
	importStaff(t, m, staffmanager.StaffMember{UserId: "u1", Name: "Ana Rojas", Rut: "1-9", IsActive: true})
	importStaff(t, m, staffmanager.StaffMember{UserId: "u2", Name: "Luis Fuentes", Rut: "2-7", IsActive: true})

	all, err := m.ListStaff(testTenant, "", "", 10)
	if err != nil || len(all) != 2 {
		t.Fatalf("ListStaff without search = %d rows err=%v", len(all), err)
	}
	found, err := m.ListStaff(testTenant, "Ana", "", 10)
	if err != nil || len(found) != 1 || found[0].Name != "Ana Rojas" {
		t.Fatalf("ListStaff(\"Ana\") = %+v err=%v", found, err)
	}
}
