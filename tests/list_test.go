package tests

import (
	"testing"

	staffmanager "github.com/veltylabs/staff_manager"
	"webtyp.com/router/loopback"
	"webtyp.com/view"
)

func TestStaff_ListAcceptsListArgs(t *testing.T) {
	m := setup(t, nil)

	// Seed rows
	rows := []staffmanager.StaffMember{
		{UserId: "u1", Rut: "11111111-1", Name: "Alice", IsActive: true},
		{UserId: "u2", Rut: "22222222-2", Name: "Bob", IsActive: true},
		{UserId: "u3", Rut: "33333333-3", Name: "Charlie", IsActive: true},
	}
	var ids []string
	for _, r := range rows {
		created, err := m.UpsertStaff(r)
		if err != nil {
			t.Fatalf("UpsertStaff: %v", err)
		}
		ids = append(ids, created.Id)
	}

	caller := loopback.ActingAs(testTenant, "user-100", m)

	// Test limit 1
	var list1 staffmanager.StaffMemberList
	var err error
	var cb1 = func(e error) { err = e }
	caller.Call(staffmanager.ModelName+"."+staffmanager.OpListStaff, &view.ListArgs{Limit: 1}, &list1, cb1)
	if err != nil {
		t.Fatalf("List with limit failed: %v", err)
	}
	if len(list1) != 1 {
		t.Errorf("expected length 1, got %d", len(list1))
	}
	if list1[0].Id != ids[0] {
		t.Errorf("expected first row id %s, got %s", ids[0], list1[0].Id)
	}

	// Test after
	var list2 staffmanager.StaffMemberList
	var err2 error
	var cb2 = func(e error) { err2 = e }
	caller.Call(staffmanager.ModelName+"."+staffmanager.OpListStaff, &view.ListArgs{After: ids[0], Limit: 10}, &list2, cb2)
	if err2 != nil {
		t.Fatalf("List with after failed: %v", err2)
	}
	if len(list2) != 2 {
		t.Errorf("expected length 2, got %d", len(list2))
	}
	if list2[0].Id == ids[0] {
		t.Errorf("did not expect first row id %s to be present", ids[0])
	}

	// Test search name
	var list3 staffmanager.StaffMemberList
	var err3 error
	var cb3 = func(e error) { err3 = e }
	caller.Call(staffmanager.ModelName+"."+staffmanager.OpListStaff, &view.ListArgs{Search: "Bob", Limit: 10}, &list3, cb3)
	if err3 != nil {
		t.Fatalf("List with search failed: %v", err3)
	}
	if len(list3) != 1 {
		t.Errorf("expected length 1, got %d", len(list3))
	}
	if list3[0].Name != "Bob" {
		t.Errorf("expected name Bob, got %s", list3[0].Name)
	}

	// Test search rut
	var list4 staffmanager.StaffMemberList
	var err4 error
	var cb4 = func(e error) { err4 = e }
	caller.Call(staffmanager.ModelName+"."+staffmanager.OpListStaff, &view.ListArgs{Search: "33333333-3", Limit: 10}, &list4, cb4)
	if err4 != nil {
		t.Fatalf("List with search rut failed: %v", err4)
	}
	if len(list4) != 1 {
		t.Errorf("expected length 1, got %d", len(list4))
	}
	if list4[0].Name != "Charlie" {
		t.Errorf("expected name Charlie, got %s", list4[0].Name)
	}
}
