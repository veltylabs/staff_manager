package tests
import (
    "testing"
    "webtyp.com/model"
    "webtyp.com/crud"
    "webtyp.com/router"
    "github.com/veltylabs/staff_manager"
)
func TestStaff_DuplicateRutRejected(t *testing.T) {
    m := setup(t, nil)
    caller := crud.LoopbackActingAs(nil, "u1", m)
    caller.Call(router.Operation{Module: staffmanager.ModelName, Entity: "staff_member", Action: model.Create}, &crud.CreateArgs{
        Records: []model.Model{&staffmanager.StaffMember{Id: "s1", Rut: "11111111-1", Name: "Bob", IsActive: true}},
    }, nil, func(err error) {
        if err != nil { t.Errorf("setup err: %v", err) }
    })
    caller.Call(router.Operation{Module: staffmanager.ModelName, Entity: "staff_member", Action: model.Create}, &crud.CreateArgs{
        Records: []model.Model{&staffmanager.StaffMember{Id: "s2", Rut: "11111111-1", Name: "Alice", IsActive: true}},
    }, nil, func(err error) {
        if err == nil { t.Errorf("expected error, got nil") }
    })
}
