package staffmanager
import (
	"webtyp.com/crud"
	"webtyp.com/model"
	"webtyp.com/storage"
)
func (m *Module) RegisterFirstAdmin(userID, code, name, ip string) error {
	if m.devices == nil { return depErr("Devices") }
	if m.registrar == nil { return depErr("Registrar") }
	return m.store.Atomic(func(tx crud.Tx) error {
		member := StaffMember{TenantId: m.tenantID, UserId: userID, Rut: code, Name: name, IsActive: true}
		member.Id = m.ids.NewID()
		if err := checkStaffChange(tx, m, nil, &member, nil); err != nil { return err }
		if err := tx.Create(&member); err != nil { return err }
		ips := []string{ip}
		if ip == "127.0.0.1" { ips = append(ips, "::1") } else if ip == "::1" { ips = append(ips, "127.0.0.1") }
		for _, addr := range ips {
			deviceID, ok := m.devices.FindByIP(addr)
			if !ok {
				var err error
				if deviceID, err = m.registrar.RegisterComputer(name+" (admin)", addr); err != nil { return err }
			}
			var sd StaffDevice
			err := m.StoreStaffDevice().ReadAll(func() model.Model { return &StaffDevice{} }, crud.Query{Conditions: []storage.Condition{storage.Eq("staff_id", member.Id), storage.Eq("device_id", deviceID)}, Limit: 1}, func(mod model.Model) { sd = *(mod.(*StaffDevice)) })
			if err == nil && sd.StaffId == "" {
				sd = StaffDevice{StaffId: member.Id, DeviceId: deviceID}
				if err := m.StoreStaffDevice().Create(&sd); err != nil { return err }
			} else if err != nil { return err }
		}
		return nil
	})
}
