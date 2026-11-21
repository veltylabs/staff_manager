package staffmanager

import (
	"webtyp.com/crud"
	"webtyp.com/model"
	"webtyp.com/orm"
	"webtyp.com/storage"
)

func (m *Module) RegisterFirstAdmin(userID, code, name, ip string) error {
	if m.devices == nil {
		return depErr("Devices")
	}
	if m.registrar == nil {
		return depErr("Registrar")
	}
	// Una sola transacción para el profesional (tabla con tenant) y su vínculo con el equipo (StaffDevice, sin tenant):
	// db.Tx da los dos almacenes sobre la misma transacción. m.store.Atomic solo da el de tenant, y escribir el vínculo
	// con m.deviceStore fuera de ella fallaba la llave foránea (el profesional aún no estaba confirmado). H-ACCESS-01.
	return m.db.Tx(func(db *orm.DB) error {
		tx := orm.TenantStore(db, m.tenantID, m.ids)
		links := orm.GlobalStore(db, m.ids)
		member := StaffMember{TenantId: m.tenantID, UserId: userID, Rut: code, Name: name, IsActive: true}
		member.Id = m.ids.NewID()
		if err := checkStaffChange(tx, m, nil, &member, nil); err != nil {
			return err
		}
		if err := tx.Create(&member); err != nil {
			return err
		}
		ips := []string{ip}
		switch ip {
		case "127.0.0.1":
			ips = append(ips, "::1")
		case "::1":
			ips = append(ips, "127.0.0.1")
		}
		for _, addr := range ips {
			deviceID, ok := m.devices.FindByIP(addr)
			if !ok {
				var err error
				if deviceID, err = m.registrar.RegisterComputer(name+" (admin)", addr); err != nil {
					return err
				}
			}
			var sd StaffDevice
			err := links.ReadAll(func() model.Model { return &StaffDevice{} }, crud.Query{Conditions: []storage.Condition{storage.Eq("staff_id", member.Id), storage.Eq("device_id", deviceID)}, Limit: 1}, func(mod model.Model) { sd = *(mod.(*StaffDevice)) })
			if err == nil && sd.StaffId == "" {
				sd = StaffDevice{StaffId: member.Id, DeviceId: deviceID}
				if err := links.Create(&sd); err != nil {
					return err
				}
			} else if err != nil {
				return err
			}
		}
		return nil
	})
}
