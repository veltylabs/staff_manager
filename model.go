package staffmanager

import (
	"webtyp.com/fmt"
	"webtyp.com/input"
	"webtyp.com/model"
)

var StaffMemberModel = model.Definition{
	Name: "staff_member",
	Fields: model.Fields{
		{Name: "id", Type: model.Text(), DB: &model.FieldDB{PK: true}, OmitEmpty: true},
		{Name: "tenant_id", Type: model.Text(), NotNull: true},
		{Name: "user_id", Type: model.Text(), NotNull: true},
		{Name: "rut", Type: input.Rut(), NotNull: true, Permitted: model.Permitted{Minimum: 1, Maximum: 12}},
		{Name: "name", Type: input.Text(), NotNull: true, Permitted: model.Permitted{Minimum: 1, Maximum: 255}},
		// Specialty y role son valores de texto libre sin enums ni slugs. Role es
		// puramente descriptivo (ej. título de trabajo) y nunca una entrada de autorización.
		// Ambos campos son opcionales para que los registros existentes sigan siendo válidos.
		{Name: "specialty", Type: input.Text(), Permitted: model.Permitted{Maximum: 120}},
		{Name: "role", Type: input.Text(), Permitted: model.Permitted{Maximum: 60}},
		{Name: "is_active", Type: input.Checkbox()},
		{Name: "updated_at", Type: model.Int(), OmitEmpty: true},
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
