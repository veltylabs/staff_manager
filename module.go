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

// deviceProbe is a minimal read model for the device table, used only as
// fallback in IsTrustedIP when DeviceReader is unavailable or fails. It
// mirrors the device table's columns needed for the check.
type deviceProbe struct {
	Id       string
	TenantId string
	Ip       string
	IsActive bool
}

func (m *deviceProbe) ModelName() string { return "device" }
func (m *deviceProbe) Schema() []model.Field {
	return []model.Field{
		{Name: "id", Type: model.Text()},
		{Name: "tenant_id", Type: model.Text()},
		{Name: "ip", Type: model.Text()},
		{Name: "is_active", Type: BaseBool_FieldBool},
	}
}
func (m *deviceProbe) Pointers() []any { return []any{&m.Id, &m.TenantId, &m.Ip, &m.IsActive} }
func (m *deviceProbe) IsNil() bool     { return m == nil }
func (m *deviceProbe) EncodeFields(w model.FieldWriter) {
	w.String("id", m.Id)
	w.String("tenant_id", m.TenantId)
	w.String("ip", m.Ip)
	w.Bool("is_active", m.IsActive)
}
func (m *deviceProbe) DecodeFields(r model.FieldReader) {
	if v, ok := r.String("id"); ok {
		m.Id = v
	}
	if v, ok := r.String("tenant_id"); ok {
		m.TenantId = v
	}
	if v, ok := r.String("ip"); ok {
		m.Ip = v
	}
	if v, ok := r.Bool("is_active"); ok {
		m.IsActive = v
	}
}
func (m *deviceProbe) Validate(action byte) error { return nil }

// IsTrustedIP implements auth.TrustedIPStore
func (m *Module) IsTrustedIP(userID, ip string) bool {
	fmt.Printf("IsTrustedIP start userID=%q ip=%q tenant=%q\n", userID, ip, m.tenantID)
	if m.db == nil || userID == "" || ip == "" {
		fmt.Printf("IsTrustedIP early false\n")
		return false
	}
	var staff StaffMember
	qb := m.db.Query(&staff).Where("user_id").Eq(userID).Where("tenant_id").Eq(m.tenantID)
	err := qb.ReadOne()
	fmt.Printf("IsTrustedIP staff err=%v staff=%+v isActive=%v\n", err, staff, staff.IsActive)
	if err != nil || !staff.IsActive {
		fmt.Printf("IsTrustedIP staff not found or inactive\n")
		return false
	}
	if m.devices != nil {
		fmt.Printf("IsTrustedIP trying DeviceReader FindByIP ip=%q\n", ip)
		if deviceID, ok := m.devices.FindByIP(ip); ok {
			fmt.Printf("IsTrustedIP FindByIP ok deviceID=%q\n", deviceID)
			var sd StaffDevice
			qb2 := m.db.Query(&sd).Where("staff_id").Eq(staff.Id).Where("device_id").Eq(deviceID)
			err2 := qb2.ReadOne()
			fmt.Printf("IsTrustedIP staff_device lookup err=%v sd=%+v\n", err2, sd)
			if err2 == nil {
				fmt.Printf("IsTrustedIP success via DeviceReader\n")
				return true
			}
		} else {
			fmt.Printf("IsTrustedIP FindByIP not ok\n")
		}
	} else {
		fmt.Printf("IsTrustedIP no devices reader\n")
	}
	// Fallback: scan assigned devices and compare IP via direct device query
	var list StaffDeviceList
	if err := m.db.Query(&StaffDevice{}).Where("staff_id").Eq(staff.Id).ReadAll(
		func() model.Model { return &StaffDevice{} },
		func(mm model.Model) { list = append(list, mm.(*StaffDevice)) },
	); err != nil {
		fmt.Printf("IsTrustedIP fallback ReadAll err=%v\n", err)
		return false
	}
	fmt.Printf("IsTrustedIP fallback list len=%d\n", len(list))
	for _, sd := range list {
		fmt.Printf("IsTrustedIP fallback check sd=%+v\n", sd)
		if m.devices != nil {
			if did, ok := m.devices.FindByIP(ip); ok && did == sd.DeviceId {
				fmt.Printf("IsTrustedIP fallback success via DeviceReader did=%q\n", did)
				return true
			}
		}
		var dev deviceProbe
		if err := m.db.Query(&dev).Where("id").Eq(sd.DeviceId).Where("tenant_id").Eq(m.tenantID).ReadOne(); err == nil {
			fmt.Printf("IsTrustedIP deviceProbe dev=%+v ip equal=%v isActive=%v\n", dev, dev.Ip == ip, dev.IsActive)
			if dev.Ip == ip && dev.IsActive {
				fmt.Printf("IsTrustedIP success via deviceProbe\n")
				return true
			}
		} else {
			fmt.Printf("IsTrustedIP deviceProbe err=%v\n", err)
		}
	}
	fmt.Printf("IsTrustedIP final false\n")
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
