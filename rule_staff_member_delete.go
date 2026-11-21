package staffmanager

import (
	"webtyp.com/crud"
	"webtyp.com/model"
	"webtyp.com/storage"
)

func (p *StaffMember) Delete(a crud.CurrentUser, tx crud.Tx, m *Module) (crud.RuleAnswer, error) {
	var toDelete []StaffDevice
	err := tx.ReadAll(
		func() model.Model { return &StaffDevice{} },
		crud.Query{Conditions: []storage.Condition{storage.Eq(StaffDevice_.StaffId, p.Id)}},
		func(mod model.Model) { toDelete = append(toDelete, *(mod.(*StaffDevice))) },
	)
	if err != nil {
		return crud.RuleAnswer{}, err
	}
	for _, sd := range toDelete {
		if err := tx.Delete(&sd, sd.StaffId+"\x00"+sd.DeviceId); err != nil {
			return crud.RuleAnswer{}, err
		}
	}
	return crud.ApplyAndPublish(), nil
}
