package staffmanager

import (
	"webtyp.com/events"
	"webtyp.com/fmt"
	"webtyp.com/model"
	"webtyp.com/orm"
	"webtyp.com/time"
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

// IsTrustedIP implements auth.TrustedIPStore. The single question login
// asks: is ip a device assigned to userID? Answered through the injected
// DeviceReader only — never a direct read of device_manager's own table
// (see Deps.Devices' doc comment; the module whitelist forbids importing a
// sibling module's storage shape).
func (m *Module) IsTrustedIP(userID, ip string) bool {
	if m.db == nil || userID == "" || ip == "" {
		return false
	}
	var staff StaffMember
	qb := m.db.Query(&staff).Where("user_id").Eq(userID).Where("tenant_id").Eq(m.tenantID).Where("is_active").Eq(true)
	if err := qb.ReadOne(); err != nil {
		return false
	}
	if m.devices == nil {
		return false
	}
	deviceID, ok := m.devices.FindByIP(ip)
	if !ok {
		return false
	}
	var sd StaffDevice
	return m.db.Query(&sd).Where("staff_id").Eq(staff.Id).Where("device_id").Eq(deviceID).ReadOne() == nil
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
	member.UpdatedAt = time.Now()

	// If Id provided, try update
	if member.Id != "" {
		var existing StaffMember
		qb := m.db.Query(&existing).Where("id").Eq(member.Id).Where("tenant_id").Eq(member.TenantId)
		err := qb.ReadOne()
		if err == nil {
			// Update
			existing.Name = member.Name
			existing.Rut = member.Rut
			existing.Specialty = member.Specialty
			existing.Role = member.Role
			existing.IsActive = member.IsActive
			existing.UserId = member.UserId
			existing.UpdatedAt = member.UpdatedAt
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
		existing.Specialty = member.Specialty
		existing.Role = member.Role
		existing.IsActive = member.IsActive
		existing.UpdatedAt = member.UpdatedAt
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

// StaffExists reports whether a staff member with this id belongs to this
// tenant. It satisfies the narrow StaffReader port that a scheduling module
// declares on its own side (StaffExists(tenantId, staffId) (bool, error)) —
// structurally, with no adapter and no import in either direction.
//
// A missing row is (false, nil), NOT an error: "this id is not one of ours" is
// the answer the caller asked for. Only a real storage failure returns a
// non-nil error, so a caller can never mistake a dead database for a clean
// "no".
func (m *Module) StaffExists(tenantID, staffID string) (bool, error) {
	if tenantID == "" || staffID == "" {
		return false, nil
	}
	_, err := m.GetStaff(tenantID, staffID)
	if err != nil {
		if err == ErrNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
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
