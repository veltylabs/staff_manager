package staffmanager
import "webtyp.com/crud"
import "webtyp.com/model"
import "webtyp.com/storage"
func (p *StaffMember) Delete(a crud.CurrentUser, tx crud.Tx, m *Module) (crud.RuleAnswer, error) {
    var toDelete []StaffDevice
    tx.ReadAll(func() model.Model { return &StaffDevice{} }, crud.Query{Conditions: []storage.Condition{storage.Eq("staff_id", p.Id)}}, func(mod model.Model) { toDelete = append(toDelete, *(mod.(*StaffDevice))) })
    for _, sd := range toDelete {
        tx.Delete(&sd, sd.StaffId + "\x00" + sd.DeviceId) // Assumes serialized composite key for crud
    }
    return crud.ApplyAndPublish, nil
}
