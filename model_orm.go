// Code generated manually — same pattern as device_manager/model_orm.go
package staffmanager

import (
	"webtyp.com/model"
	"webtyp.com/orm"
)

type StaffMember struct {
	Id        string
	TenantId  string
	UserId    string
	Rut       string
	Name      string
	IsActive  bool
	UpdatedAt int64
}

func (m *StaffMember) ModelName() string { return "staff_member" }
func (m *StaffMember) Schema() []model.Field { return StaffMemberModel.Fields }
func (m *StaffMember) Pointers() []any {
	return []any{&m.Id, &m.TenantId, &m.UserId, &m.Rut, &m.Name, &m.IsActive, &m.UpdatedAt}
}
func (m *StaffMember) IsNil() bool { return m == nil }
func (m *StaffMember) EncodeFields(w model.FieldWriter) {
	if m.Id != "" {
		w.String("id", m.Id)
	}
	w.String("tenant_id", m.TenantId)
	w.String("user_id", m.UserId)
	w.String("rut", m.Rut)
	w.String("name", m.Name)
	w.Bool("is_active", m.IsActive)
	if m.UpdatedAt != 0 {
		w.Int("updated_at", m.UpdatedAt)
	}
}
func (m *StaffMember) DecodeFields(r model.FieldReader) {
	if v, ok := r.String("id"); ok {
		m.Id = v
	}
	if v, ok := r.String("tenant_id"); ok {
		m.TenantId = v
	}
	if v, ok := r.String("user_id"); ok {
		m.UserId = v
	}
	if v, ok := r.String("rut"); ok {
		m.Rut = v
	}
	if v, ok := r.String("name"); ok {
		m.Name = v
	}
	if v, ok := r.Bool("is_active"); ok {
		m.IsActive = v
	}
	if v, ok := r.Int("updated_at"); ok {
		m.UpdatedAt = v
	}
}
func (m *StaffMember) Validate(action byte) error { return model.ValidateFields(action, m) }

var StaffMember_ = struct {
	Id        string
	TenantId  string
	UserId    string
	Rut       string
	Name      string
	IsActive  string
	UpdatedAt string
}{
	Id: "id", TenantId: "tenant_id", UserId: "user_id", Rut: "rut", Name: "name", IsActive: "is_active", UpdatedAt: "updated_at",
}

func ReadOneStaffMember(qb *orm.QB, model *StaffMember) (*StaffMember, error) {
	err := qb.ReadOne()
	if err != nil {
		return nil, err
	}
	return model, nil
}

type StaffMemberList []*StaffMember

func (s *StaffMemberList) Len() int                     { return len(*s) }
func (s *StaffMemberList) At(i int) model.Fielder      { return (*s)[i] }
func (s *StaffMemberList) Append() model.Fielder       { v := &StaffMember{}; *s = append(*s, v); return v }
func (s *StaffMemberList) IsNil() bool                 { return s == nil }
func (s *StaffMemberList) EncodeFields(_ model.FieldWriter) {}
func (s *StaffMemberList) DecodeFields(_ model.FieldReader) {}

type StaffDevice struct {
	StaffId  string
	DeviceId string
}

func (m *StaffDevice) ModelName() string { return "staff_device" }
func (m *StaffDevice) Schema() []model.Field { return StaffDeviceModel.Fields }
func (m *StaffDevice) Pointers() []any { return []any{&m.StaffId, &m.DeviceId} }
func (m *StaffDevice) IsNil() bool { return m == nil }
func (m *StaffDevice) EncodeFields(w model.FieldWriter) {
	w.String("staff_id", m.StaffId)
	w.String("device_id", m.DeviceId)
}
func (m *StaffDevice) DecodeFields(r model.FieldReader) {
	if v, ok := r.String("staff_id"); ok {
		m.StaffId = v
	}
	if v, ok := r.String("device_id"); ok {
		m.DeviceId = v
	}
}
func (m *StaffDevice) Validate(action byte) error { return model.ValidateFields(action, m) }

var StaffDevice_ = struct {
	StaffId  string
	DeviceId string
}{
	StaffId: "staff_id", DeviceId: "device_id",
}

func ReadOneStaffDevice(qb *orm.QB, model *StaffDevice) (*StaffDevice, error) {
	err := qb.ReadOne()
	if err != nil {
		return nil, err
	}
	return model, nil
}

type StaffDeviceList []*StaffDevice

func (s *StaffDeviceList) Len() int                     { return len(*s) }
func (s *StaffDeviceList) At(i int) model.Fielder      { return (*s)[i] }
func (s *StaffDeviceList) Append() model.Fielder       { v := &StaffDevice{}; *s = append(*s, v); return v }
func (s *StaffDeviceList) IsNil() bool                 { return s == nil }
func (s *StaffDeviceList) EncodeFields(_ model.FieldWriter) {}
func (s *StaffDeviceList) DecodeFields(_ model.FieldReader) {}

type ListStaffArgs struct {
	TenantId string
}

func (m *ListStaffArgs) ModelName() string { return "list_staff_args" }
func (m *ListStaffArgs) Schema() []model.Field { return ListStaffArgsModel.Fields }
func (m *ListStaffArgs) Pointers() []any { return []any{&m.TenantId} }
func (m *ListStaffArgs) IsNil() bool { return m == nil }
func (m *ListStaffArgs) EncodeFields(w model.FieldWriter) { w.String("tenant_id", m.TenantId) }
func (m *ListStaffArgs) DecodeFields(r model.FieldReader) {
	if v, ok := r.String("tenant_id"); ok {
		m.TenantId = v
	}
}
func (m *ListStaffArgs) Validate(action byte) error { return nil }

type GetStaffArgs struct {
	TenantId string
	Id       string
}

func (m *GetStaffArgs) ModelName() string { return "get_staff_args" }
func (m *GetStaffArgs) Schema() []model.Field { return GetStaffArgsModel.Fields }
func (m *GetStaffArgs) Pointers() []any { return []any{&m.TenantId, &m.Id} }
func (m *GetStaffArgs) IsNil() bool { return m == nil }
func (m *GetStaffArgs) EncodeFields(w model.FieldWriter) {
	w.String("tenant_id", m.TenantId)
	w.String("id", m.Id)
}
func (m *GetStaffArgs) DecodeFields(r model.FieldReader) {
	if v, ok := r.String("tenant_id"); ok {
		m.TenantId = v
	}
	if v, ok := r.String("id"); ok {
		m.Id = v
	}
}
func (m *GetStaffArgs) Validate(action byte) error { return nil }
