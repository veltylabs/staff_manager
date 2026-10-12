package tests

import (
	"testing"

	"webtyp.com/fmt"
	"webtyp.com/orm"
	"webtyp.com/storage"
	"webtyp.com/storage/mem"
	staffmanager "github.com/veltylabs/staff_manager"
)

func TestStaffCard_EmptyArgsNotFound(t *testing.T) {
	m := setup(t, nil)
	// Tenant vacío
	_, _, _, err := m.StaffCard("", "id-123")
	if err != staffmanager.ErrNotFound {
		t.Errorf("esperaba ErrNotFound con tenant vacío, obtuve %v", err)
	}
	// ID vacío
	_, _, _, err = m.StaffCard(testTenant, "")
	if err != staffmanager.ErrNotFound {
		t.Errorf("esperaba ErrNotFound con id vacío, obtuve %v", err)
	}
}

func TestStaffCard_ReturnsNameSpecialtyActive(t *testing.T) {
	m := setup(t, nil)

	// Crear uno activo
	s1, err := m.UpsertStaff(staffmanager.StaffMember{
		Name:      "Ana",
		Specialty: "Medicina General",
		IsActive:  true,
		UserId:    "usr1",
		Rut:       "1-9",
	})
	if err != nil {
		t.Fatalf("UpsertStaff: %v", err)
	}

	name, specialty, active, err := m.StaffCard(testTenant, s1.Id)
	if err != nil {
		t.Fatalf("StaffCard s1 falló: %v", err)
	}
	if name != "Ana" || specialty != "Medicina General" || active != true {
		t.Errorf("StaffCard s1 incorrecto. got: %q %q %v", name, specialty, active)
	}

	// Crear uno inactivo
	s2, err := m.UpsertStaff(staffmanager.StaffMember{
		Name:      "Luis",
		Specialty: "Pediatría",
		IsActive:  false,
		UserId:    "usr2",
		Rut:       "2-7",
	})
	if err != nil {
		t.Fatalf("UpsertStaff: %v", err)
	}

	name, specialty, active, err = m.StaffCard(testTenant, s2.Id)
	if err != nil {
		t.Fatalf("StaffCard s2 falló: %v", err)
	}
	if name != "Luis" || specialty != "Pediatría" || active != false {
		t.Errorf("StaffCard s2 incorrecto. got: %q %q %v", name, specialty, active)
	}
}

func TestStaffCard_UnknownOrOtherTenantNotFound(t *testing.T) {
	m := setup(t, nil)

	s1, err := m.UpsertStaff(staffmanager.StaffMember{
		Name:      "Ana",
		IsActive:  true,
		UserId:    "usr1",
		Rut:       "1-9",
	})
	if err != nil {
		t.Fatalf("UpsertStaff: %v", err)
	}

	// Otro tenant
	_, _, _, err = m.StaffCard("otro-tenant", s1.Id)
	if err != staffmanager.ErrNotFound {
		t.Errorf("esperaba ErrNotFound con otro tenant, obtuve %v", err)
	}

	// ID inexistente
	_, _, _, err = m.StaffCard(testTenant, "no-existe")
	if err != staffmanager.ErrNotFound {
		t.Errorf("esperaba ErrNotFound con id desconocido, obtuve %v", err)
	}
}

// fakeStorageConn envuelve un storage.Conn real pero fuerza error en Query
type fakeStorageConn struct {
	storage.Conn
}

type brokenScanner struct{}
func (brokenScanner) Scan(dest ...any) error { return fmt.Err("error simulado") }

func (f *fakeStorageConn) QueryRow(query string, args ...any) storage.Scanner {
	return brokenScanner{}
}

func TestStaffCard_StorageErrorPropagates(t *testing.T) {
	// Crear un db que falla en Query
	brokenDB := orm.New(&fakeStorageConn{Conn: mem.New()})
	m, err := staffmanager.New(brokenDB, staffmanager.Deps{
		IDs:         &mockIDGen{},
		Publisher:   &mockPublisher{},
		TenantID:    testTenant,
		ValidateRUT: identityValidateRUT,
		Devices:     nil,
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	_, _, _, err = m.StaffCard(testTenant, "id-123")
	if err == nil || err.Error() != "error simulado" {
		t.Errorf("esperaba error simulado, obtuve %v", err)
	}
}