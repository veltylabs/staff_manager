package staffmanager

import (
	"webtyp.com/events"
	"webtyp.com/fmt"
	"webtyp.com/model"
	"webtyp.com/orm"
)

type DeviceReader interface {
	FindByIP(ip string) (string, bool)
}

type Deps struct {
	IDs         model.IDGenerator
	Publisher   events.Publisher
	TenantID    string
	ValidateRUT func(string) (string, error)
	Devices     DeviceReader
}

type Module struct {
	db          *orm.DB
	ids         model.IDGenerator
	pub         events.Publisher
	tenantID    string
	validateRUT func(string) (string, error)
	devices     DeviceReader
}

func New(db *orm.DB, deps Deps) (*Module, error) {
	if deps.IDs == nil {
		return nil, fmt.Err("staff_manager: Deps.IDs is required")
	}
	if deps.TenantID == "" {
		return nil, fmt.Err("staff_manager: Deps.TenantID is required")
	}
	if deps.ValidateRUT == nil {
		return nil, fmt.Err("staff_manager: Deps.ValidateRUT is required")
	}
	return &Module{
		db:          db,
		ids:         deps.IDs,
		pub:         deps.Publisher,
		tenantID:    deps.TenantID,
		validateRUT: deps.ValidateRUT,
		devices:     deps.Devices,
	}, nil
}

// IsTrustedIP implements auth.TrustedIPStore
func (m *Module) IsTrustedIP(userID, ip string) bool {
	if m.db == nil || userID == "" || ip == "" {
		return false
	}
	// Find staff member for user
	var staff StaffMember
	qb := m.db.Query(&staff).Where("user_id").Eq(userID).Where("tenant_id").Eq(m.tenantID).Where("is_active").Eq(true)
	err := qb.ReadOne()
	if err != nil {
		return false
	}
	// Find device id for IP via DeviceReader if available, else check via device table directly?
	// Prefer DeviceReader if provided, but fallback to checking device id via IP lookup in staff_device join:
	// If Devices is available, resolve IP -> deviceID, then check staff_device
	if m.devices != nil {
		deviceID, ok := m.devices.FindByIP(ip)
		if !ok {
			return false
		}
		var sd StaffDevice
		qb2 := m.db.Query(&sd).Where("staff_id").Eq(staff.Id).Where("device_id").Eq(deviceID)
		err = qb2.ReadOne()
		return err == nil
	}
	// Fallback: check if any device assigned to this staff has this IP by joining? Not available without device table.
	// Return false if no Devices reader.
	return false
}

func (m *Module) UpsertStaff(member StaffMember) (StaffMember, error) {
	if member.TenantId == "" {
		member.TenantId = m.tenantID
	}
	// Validate RUT
	normalized, err := m.validateRUT(member.Rut)
	if err != nil {
		return StaffMember{}, err
	}
	member.Rut = normalized

	// If Id provided, try update
	if member.Id != "" {
		var existing StaffMember
		qb := m.db.Query(&existing).Where("id").Eq(member.Id).Where("tenant_id").Eq(member.TenantId)
		err := qb.ReadOne()
		if err == nil {
			// Update
			existing.Name = member.Name
			existing.Rut = member.Rut
			existing.IsActive = member.IsActive
			existing.UserId = member.UserId
			if err := m.db.Update(&existing, orm.Eq("id", existing.Id), orm.Eq("tenant_id", existing.TenantId)); err != nil {
				return StaffMember{}, err
			}
			return existing, nil
		}
		if err != orm.ErrNotFound {
			return StaffMember{}, err
		}
	}
	// Check existing by tenant+rut or tenant+user_id
	var existing StaffMember
	qb := m.db.Query(&existing).Where("tenant_id").Eq(member.TenantId).Where("rut").Eq(member.Rut)
	err = qb.ReadOne()
	if err == nil {
		// Update existing
		existing.Name = member.Name
		existing.IsActive = member.IsActive
		if member.UserId != "" {
			existing.UserId = member.UserId
		}
		if err := m.db.Update(&existing, orm.Eq("id", existing.Id), orm.Eq("tenant_id", existing.TenantId)); err != nil {
			return StaffMember{}, err
		}
		return existing, nil
	}
	if err != orm.ErrNotFound && err != nil {
		return StaffMember{}, err
	}
	// Create new
	if member.Id == "" {
		member.Id = m.ids.NewID()
	}
	if err := m.db.Create(&member); err != nil {
		return StaffMember{}, err
	}
	return member, nil
}

func (m *Module) AssignDevice(staffID, deviceID string) error {
	if staffID == "" || deviceID == "" {
		return fmt.Err("staff_manager: staffID and deviceID required")
	}
	// Check already exists
	var sd StaffDevice
	qb := m.db.Query(&sd).Where("staff_id").Eq(staffID).Where("device_id").Eq(deviceID)
	err := qb.ReadOne()
	if err == nil {
		return nil // already assigned
	}
	if err != orm.ErrNotFound && err != nil {
		return err
	}
	sd = StaffDevice{StaffId: staffID, DeviceId: deviceID}
	return m.db.Create(&sd)
}

func (m *Module) GetStaff(tenantID, id string) (StaffMember, error) {
	var s StaffMember
	qb := m.db.Query(&s).Where("id").Eq(id).Where("tenant_id").Eq(tenantID)
	err := qb.ReadOne()
	if err != nil {
		if err == orm.ErrNotFound {
			return StaffMember{}, ErrNotFound
		}
		return StaffMember{}, err
	}
	return s, nil
}

func (m *Module) ListStaff(tenantID string) ([]StaffMember, error) {
	qb := m.db.Query(&StaffMember{}).Where("tenant_id").Eq(tenantID)
	results, err := ReadAllStaffMember(qb)
	if err != nil {
		return nil, err
	}
	out := make([]StaffMember, len(results))
	for i, v := range results {
		out[i] = *v
	}
	return out, nil
}

func ReadAllStaffMember(qb *orm.QB) (StaffMemberList, error) {
	var results StaffMemberList
	err := qb.ReadAll(
		func() model.Model { return &StaffMember{} },
		func(m model.Model) { results = append(results, m.(*StaffMember)) },
	)
	return results, err
}
