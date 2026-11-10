package staffmanager
import (
	"strings"
	"webtyp.com/events"
	"webtyp.com/fmt"
	"webtyp.com/model"
	"webtyp.com/orm"
	"webtyp.com/crud"
	"webtyp.com/storage"
)
type Deps struct {
	IDs         model.IDGenerator
	Publisher   events.Publisher
	TenantID    string
	ValidateRUT func(string) (string, error)
	Devices     DeviceReader
	Registrar   DeviceRegistrar
}
type Module struct {
	db          *orm.DB
	store       *orm.Store
	ids         model.IDGenerator
	tenantID    string
	validateRUT func(string) (string, error)
	devices     DeviceReader
	registrar   DeviceRegistrar
}
func depErr(name string) error {
	return fmt.Err("staff_manager: Deps."+name, "is required")
}
func New(db *orm.DB, deps Deps) (*Module, error) {
	if deps.IDs == nil { return nil, depErr("IDs") }
	if deps.TenantID == "" { return nil, depErr("TenantID") }
	if deps.ValidateRUT == nil { return nil, depErr("ValidateRUT") }
	return &Module{
		db:          db,
		store:       orm.TenantStore(db, deps.TenantID, deps.IDs),
		ids:         deps.IDs,
		tenantID:    deps.TenantID,
		validateRUT: deps.ValidateRUT,
		devices:     deps.Devices,
		registrar:   deps.Registrar,
	}, nil
}
func (m *Module) IsTrustedIP(userID, ip string) bool {
	if m.db == nil || userID == "" || ip == "" { return false }
	var staff StaffMember
	err := m.store.ReadAll(func() model.Model { return &StaffMember{} }, crud.Query{Conditions: []storage.Condition{storage.Eq("user_id", userID), storage.Eq("is_active", true)}}, func(mod model.Model) { staff = *(mod.(*StaffMember)) })
	if err != nil || staff.Id == "" { return false }
	if m.devices == nil { return false }
	deviceID, ok := m.devices.FindByIP(ip)
	if !ok { return false }
	var sd StaffDevice
	err = m.store.ReadAll(func() model.Model { return &StaffDevice{} }, crud.Query{Conditions: []storage.Condition{storage.Eq("staff_id", staff.Id), storage.Eq("device_id", deviceID)}}, func(mod model.Model) { sd = *(mod.(*StaffDevice)) })
	return sd.StaffId != ""
}
func (m *Module) AssignDevice(staffID, deviceID string) error {
	if staffID == "" || deviceID == "" { return fmt.Err("staff_manager: staffID and deviceID required") }
	return m.store.Atomic(func(tx crud.Tx) error {
		var sd StaffDevice
		err := tx.ReadAll(func() model.Model { return &StaffDevice{} }, crud.Query{Conditions: []storage.Condition{storage.Eq("staff_id", staffID), storage.Eq("device_id", deviceID)}}, func(mod model.Model) { sd = *(mod.(*StaffDevice)) })
		if err == nil && sd.StaffId != "" { return nil }
		sd = StaffDevice{StaffId: staffID, DeviceId: deviceID}
		return tx.Create(&sd)
	})
}
func (m *Module) GetStaff(tenantID, id string) (StaffMember, error) {
	if tenantID != m.tenantID { return StaffMember{}, ErrNotFound }
	var s StaffMember
	err := m.store.ReadOne(&s, id)
	if err != nil {
		if model.KindOf(err) == model.ErrorNotFound { return StaffMember{}, ErrNotFound }
		return StaffMember{}, err
	}
	return s, nil
}
func (m *Module) StaffCard(tenantID, staffID string) (name, specialty string, active bool, err error) {
	if tenantID != m.tenantID || staffID == "" { return "", "", false, ErrNotFound }
	s, err := m.GetStaff(tenantID, staffID)
	if err != nil { return "", "", false, err }
	return s.Name, s.Specialty, s.IsActive, nil
}
func (m *Module) StaffByUser(tenantID, userID string) (staffID, name, specialty string, found bool, err error) {
	if tenantID != m.tenantID || userID == "" { return "", "", "", false, nil }
	var rows []StaffMember
	err = m.store.ReadAll(func() model.Model { return &StaffMember{} }, crud.Query{Conditions: []storage.Condition{storage.Eq("user_id", userID)}}, func(mod model.Model) { rows = append(rows, *(mod.(*StaffMember))) })
	if err != nil {
		if model.KindOf(err) == model.ErrorNotFound { return "", "", "", false, nil }
		return "", "", "", false, err
	}
	if len(rows) == 0 { return "", "", "", false, nil }
	s := rows[0]
	return s.Id, s.Name, s.Specialty, true, nil
}
func (m *Module) StaffExists(tenantID, staffID string) (bool, error) {
	if tenantID != m.tenantID || staffID == "" { return false, nil }
	_, err := m.GetStaff(tenantID, staffID)
	if err != nil {
		if model.KindOf(err) == model.ErrorNotFound { return false, nil }
		return false, err
	}
	return true, nil
}
func (m *Module) ListStaff(tenantID, search, after string, limit int64) ([]StaffMember, error) {
	if tenantID != m.tenantID { return nil, nil }
	var conds []storage.Condition
	if search != "" {
		if strings.Contains(search, "-") || (len(search) > 0 && search[0] >= '0' && search[0] <= '9') {
			normalized, _ := m.validateRUT(search)
			if normalized != "" { conds = append(conds, storage.Eq("rut", normalized))
			} else { conds = append(conds, storage.Like("name", "%"+search+"%")) }
		} else { conds = append(conds, storage.Like("name", "%"+search+"%")) }
	}
	var out []StaffMember
	err := m.store.ReadAll(func() model.Model { return &StaffMember{} }, crud.Query{Conditions: conds, Limit: limit, After: after}, func(mod model.Model) { out = append(out, *(mod.(*StaffMember))) })
	if err != nil { return nil, err }
	return out, nil
}
