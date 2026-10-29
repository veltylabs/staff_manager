//go:build browser

package tests

import (
	"testing"

	staffmanager "github.com/veltylabs/staff_manager"
	"webtyp.com/router/loopback"
	"webtyp.com/view"
	"webtyp.com/model"
)

func TestScreen_SaveAndDelete(t *testing.T) {
	// Setup module
	m := setup(t, nil)

	caller := loopback.ActingAs(testTenant, "user-100", m)

	// Mount caller using NewRecords wrapper
	caller = view.NewRecords().Caller(caller)

	// Create view list
	v := view.NewCallerLister(caller,
		view.Ops{Module: staffmanager.ModelName, List: staffmanager.OpListStaff, Save: staffmanager.OpUpsertStaff, Delete: staffmanager.OpDeleteStaff},
		func() model.ModelSlice { return &staffmanager.StaffMemberList{} })

	// 1. Save
	staff := &staffmanager.StaffMember{
		TenantId: testTenant,
		UserId:   "user-100",
		Rut:      "12345678-5",
		Name:     "Dr. John Doe",
		IsActive: true,
	}

	saver, ok := v.(view.Saver)
	if !ok {
		t.Fatalf("Lister does not implement view.Saver")
	}

	var saveErr error
	saver.Save([]model.Model{staff}, func(err error) {
		saveErr = err
	})
	if saveErr != nil {
		t.Fatalf("Save failed: %v", saveErr)
	}
	err := saveErr

	// Wait a moment for save or get latest
	list, err := m.ListStaff(testTenant)
	if err != nil {
		t.Fatalf("ListStaff: %v", err)
	}
	if len(list) == 0 {
		t.Fatalf("No staff created")
	}
	staff.Id = list[len(list)-1].Id

	// Create a dummy staff device to test deletion cascade
	err = m.AssignDevice(staff.Id, "device-123")
	if err != nil {
		t.Fatalf("AssignDevice failed: %v", err)
	}

	// 2. Delete
	deleter, ok := v.(view.Deleter)
	if !ok {
		t.Fatalf("Lister does not implement view.Deleter")
	}

	var delErr error
	deleter.Delete([]string{staff.Id}, func(err error) {
		delErr = err
	})
	if delErr != nil {
		t.Fatalf("Delete failed: %v", delErr)
	}

	// Verify it's removed from the DB
	_, err = m.GetStaff(testTenant, staff.Id)
	if err == nil {
		t.Fatalf("Expected error after deletion, got nil")
	}

	if model.KindOf(err) != model.ErrorNotFound {
		t.Fatalf("Expected model.ErrorNotFound, got %v", err)
	}
}

// El tenant lo pone el servidor (D32): un registro que llega con otro tenant
// se guarda en el de esta instalación, y el otro no recibe nada.
func TestScreen_SaveIgnoresTenantFromRequest(t *testing.T) {
	m := setup(t, nil)
	caller := view.NewRecords().Caller(loopback.ActingAs(testTenant, "user-100", m))
	v := view.NewCallerLister(caller,
		view.Ops{Module: staffmanager.ModelName, List: staffmanager.OpListStaff, Save: staffmanager.OpUpsertStaff, Delete: staffmanager.OpDeleteStaff},
		func() model.ModelSlice { return &staffmanager.StaffMemberList{} })

	const otherTenant = "other-tenant"
	staff := &staffmanager.StaffMember{TenantId: otherTenant, UserId: "user-200", Rut: "11111111-1", Name: "Ana Torres", IsActive: true}
	var saveErr error
	v.(view.Saver).Save([]model.Model{staff}, func(err error) { saveErr = err })
	if saveErr != nil {
		t.Fatalf("Save failed: %v", saveErr)
	}

	mine, err := m.ListStaff(testTenant)
	if err != nil {
		t.Fatalf("ListStaff(%s): %v", testTenant, err)
	}
	if len(mine) != 1 {
		t.Fatalf("stored under the server's tenant: got %d rows, want 1", len(mine))
	}
	other, err := m.ListStaff(otherTenant)
	if err != nil {
		t.Fatalf("ListStaff(%s): %v", otherTenant, err)
	}
	if len(other) != 0 {
		t.Fatalf("the tenant sent by the request got %d rows, want 0", len(other))
	}
}
