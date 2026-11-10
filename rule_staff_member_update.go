package staffmanager

import "webtyp.com/crud"

func (next *StaffMember) Update(a crud.CurrentUser, tx crud.Tx, old *StaffMember, fields []string, m *Module) (crud.RuleAnswer, error) {
	if err := checkStaffChange(tx, m, old, next, fields); err != nil {
		return 0, err
	}
	return crud.ApplyAndPublish, nil
}
