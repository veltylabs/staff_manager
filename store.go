package staffmanager

import (
	"webtyp.com/crud"
)

func (m *Module) StoreStaffMember() crud.Store { return m.store }
func (m *Module) StoreStaffDevice() crud.Store { return m.deviceStore }
