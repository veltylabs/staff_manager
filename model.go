package staffmanager

import (
	"webtyp.com/fmt"
	"webtyp.com/input"
	"webtyp.com/model"
)

var (
	BaseBool_FieldBool = model.Bool()
	BaseInt_FieldInt   = model.Int()
)

var StaffMemberModel = model.Definition{
	Name: "staff_member",
	Fields: model.Fields{
		{Name: "id", Type: model.Text(), DB: &model.FieldDB{PK: true}, OmitEmpty: true},
		{Name: "tenant_id", Type: model.Text(), NotNull: true},
		{Name: "user_id", Type: model.Text(), NotNull: true},
		{Name: "rut", Type: input.Text(), NotNull: true, Permitted: model.Permitted{Minimum: 1, Maximum: 12}},
		{Name: "name", Type: input.Text(), NotNull: true, Permitted: model.Permitted{Minimum: 1, Maximum: 255}},
		{Name: "is_active", Type: input.Checkbox(), NotNull: true},
		{Name: "updated_at", Type: BaseInt_FieldInt, OmitEmpty: true},
	},
}

var StaffDeviceModel = model.Definition{
	Name: "staff_device",
	Fields: model.Fields{
		{Name: "staff_id", Type: model.Text(), DB: &model.FieldDB{PK: true}, NotNull: true},
		{Name: "device_id", Type: model.Text(), DB: &model.FieldDB{PK: true}, NotNull: true},
	},
}

var ListStaffArgsModel = model.Definition{
	Name: "list_staff_args",
	Fields: model.Fields{
		{Name: "tenant_id", Type: model.Text()},
	},
}

var GetStaffArgsModel = model.Definition{
	Name: "get_staff_args",
	Fields: model.Fields{
		{Name: "tenant_id", Type: model.Text()},
		{Name: "id", Type: model.Text()},
	},
}

var (
	ErrNotFound = fmt.Err("staff not found")
)

const (
	TopicStaffCreated = "staff_manager.staff.created"
	TopicStaffUpdated = "staff_manager.staff.updated"
)
