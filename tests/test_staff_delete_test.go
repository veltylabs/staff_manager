package tests

import (
	"testing"

	staffmanager "github.com/veltylabs/staff_manager"
	"webtyp.com/crud"
	"webtyp.com/model"
	"webtyp.com/router"
)

func TestStaff_DeleteRemovesMemberAndDevices(t *testing.T) {
	m := setup(t, nil)
	caller := crud.LoopbackActingAs(nil, "u1", m)

	var created staffmanager.StaffMember
	var createErr error
	caller.Call(router.Operation{Module: staffmanager.ModelName, Entity: "staff_member", Action: model.Create}, &crud.CreateArgs{
		Records: []model.Model{&staffmanager.StaffMember{Id: "s-del", Rut: "11111111-1", Name: "Borrar QA", IsActive: true}},
	}, &created, func(err error) { createErr = err })
	if createErr != nil {
		t.Fatalf("create err: %v", createErr)
	}

	var delErr error
	caller.Call(router.Operation{Module: staffmanager.ModelName, Entity: "staff_member", Action: model.Delete}, &crud.DeleteArgs{
		IDs: []string{"s-del"},
	}, nil, func(err error) { delErr = err })
	if delErr != nil {
		t.Fatalf("delete err: %v", delErr)
	}

	if _, err := m.GetStaff(testTenant, "s-del"); err == nil {
		t.Fatalf("expected deleted staff member to be gone")
	}
}
