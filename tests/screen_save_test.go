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
