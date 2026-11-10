package staffmanager

import (
	"webtyp.com/crud"
	"webtyp.com/model"
	"webtyp.com/storage"
	"webtyp.com/time"
)

// ImportStaffMember carga una persona del equipo con las mismas reglas de RUT que la
// operación create: si ya existe una con ese RUT, la actualiza; si no, la crea.
//
// Ejemplo: err := m.ImportStaffMember(&staffmanager.StaffMember{Rut: "17654321-3", Name: "Ana"})
func (m *Module) ImportStaffMember(row *StaffMember) error {
	row.TenantId = m.tenantID
	if row.Id == "" {
		row.Id = m.ids.NewID()
	}
	if err := model.ValidateFields(model.ActionCreate, row); err != nil {
		return err
	}
	row.UpdatedAt = time.NowNano()
	return m.store.Atomic(func(tx crud.Tx) error {
		var found []StaffMember
		err := tx.ReadAll(func() model.Model { return &StaffMember{} },
			crud.Query{Conditions: []storage.Condition{storage.Eq(StaffMember_.Rut, row.Rut)}, Limit: 1},
			func(r model.Model) { found = append(found, *(r.(*StaffMember))) })
		if err != nil {
			return err
		}
		if len(found) > 0 {
			row.Id = found[0].Id
			return tx.UpdateFields(row, row.Id, []string{StaffMember_.Name, StaffMember_.Specialty, StaffMember_.Role, StaffMember_.IsActive, StaffMember_.UserId})
		}
		return tx.Create(row)
	})
}
