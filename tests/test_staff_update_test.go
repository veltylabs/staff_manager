package tests
import (
    "testing"
    "webtyp.com/model"
    "webtyp.com/crud"
    "webtyp.com/router"
    "github.com/veltylabs/staff_manager"
)
func TestStaff_UpdateOnlyListedFields(t *testing.T) {
    m := setup(t, nil)
    caller := crud.LoopbackActingAs(nil, "u1", m)
    var created staffmanager.StaffMember
    caller.Call(router.Operation{Module: staffmanager.ModelName, Entity: "staff_member", Action: model.Create}, &crud.CreateArgs{
        Records: []model.Model{&staffmanager.StaffMember{Id: "s1", Rut: "11111111-1", Name: "Bob", IsActive: true}},
    }, &created, func(err error) {
        if err != nil { t.Errorf("setup err: %v", err) }
    })
    caller.Call(router.Operation{Module: staffmanager.ModelName, Entity: "staff_member", Action: model.Update}, &crud.UpdateArgs{
        IDs: []string{"s1"},
        Record: &staffmanager.StaffMember{Specialty: "Spec", Name: "Changed"},
        Fields: []string{"specialty"},
    }, nil, func(err error) {})
    got, _ := m.GetStaff(testTenant, "s1")
    if got.Name != "Bob" { t.Errorf("expected Bob, got %s", got.Name) }
}
