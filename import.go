package staffmanager
import "webtyp.com/crud"
import "webtyp.com/model"
import "webtyp.com/storage"
import "webtyp.com/time"
func (m *Module) ImportStaffMember(row *StaffMember) error {
    row.TenantId = m.tenantID
    model.ValidateFields(model.ActionCreate, row)
    row.UpdatedAt = time.NowNano()
    return m.store.Atomic(func(tx crud.Tx) error {
        var out []StaffMember
        err := tx.ReadAll(func() model.Model { return &StaffMember{} }, crud.Query{Conditions: []storage.Condition{storage.Eq("rut", row.Rut)}, Limit: 1}, func(mod model.Model) { out = append(out, *(mod.(*StaffMember))) })
        if err != nil { return err }
        if len(out) > 0 {
            existing := out[0]
            row.Id = existing.Id
            return tx.UpdateFields(row, existing.Id, []string{"name", "specialty", "role", "is_active", "user_id"})
        }
        row.Id = m.ids.NewID()
        return tx.Create(row)
    })
}
