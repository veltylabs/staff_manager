package staffmanager

import "webtyp.com/crud"
import "webtyp.com/model"
import "webtyp.com/storage"

func (p *StaffMember) Create(a crud.CurrentUser, tx crud.Tx, m *Module) (crud.RuleAnswer, error) {
	p.UserId = a.UserID
	p.TenantId = m.tenantID
	if err := checkStaffChange(tx, m, nil, p, nil); err != nil {
		return 0, err
	}
	return crud.ApplyAndPublish, nil
}
func checkStaffChange(tx crud.Tx, m *Module, old, next *StaffMember, fields []string) error {
	hasRut := false
	if fields == nil {
		hasRut = true
	} else {
		for _, f := range fields {
			if f == "rut" {
				hasRut = true
				break
			}
		}
	}
	if hasRut {
		normalized, err := m.validateRUT(next.Rut)
		if err != nil {
			return err
		}
		next.Rut = normalized
		var out []StaffMember
		conds := []storage.Condition{storage.Eq("rut", next.Rut)}
		if old != nil {
			conds = append(conds, storage.Neq("id", old.Id))
		}
		tx.ReadAll(func() model.Model { return &StaffMember{} }, crud.Query{Conditions: conds}, func(mod model.Model) { out = append(out, *(mod.(*StaffMember))) })
		if len(out) > 0 {
			return model.Invalid("conflict")
		}
	}
	return nil
}
